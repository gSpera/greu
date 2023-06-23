package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Execution struct {
	cmd        *exec.Cmd
	writeStdin io.WriteCloser
	readStdout io.ReadCloser
	readStderr io.ReadCloser
	tmpFile    *os.File
}

func NewExecution(cmd string, args ...string) (Execution, error) {
	e := Execution{}

	var err error

	e.tmpFile, err = os.CreateTemp("", "greu-")
	if err != nil {
		return e, fmt.Errorf("cannot create temp file: %w", err)
	}

	args = e.ReplaceMultiple(args...)
	e.cmd = exec.Command(cmd, args...)

	e.cmd.Env = append(e.cmd.Env, fmt.Sprintf("GREU_TMP=%q", e.tmpFile.Name()))

	e.writeStdin, err = e.cmd.StdinPipe()
	if err != nil {
		return e, fmt.Errorf("cannot get stdin pipe: %w", err)
	}
	e.readStdout, err = e.cmd.StdoutPipe()
	if err != nil {
		return e, fmt.Errorf("cannot get stdout pipe: %w", err)
	}
	e.readStderr, err = e.cmd.StderrPipe()
	if err != nil {
		return e, fmt.Errorf("cannot get stderr pipe: %w", err)
	}

	err = e.cmd.Start()
	if err != nil {
		return e, fmt.Errorf("cannot start command: %w", err)
	}

	return e, nil
}

func (e Execution) Write(p []byte) (int, error) {
	return e.writeStdin.Write(p)
}

func (e Execution) Read(p []byte) (int, error) {
	return e.readStdout.Read(p)
}

func (e Execution) Exit(timeout time.Duration) error {
	io.Copy(os.Stdout, e.readStdout)
	io.Copy(os.Stderr, e.readStderr)
	e.readStdout.Close()
	defer func() {
		name := e.tmpFile.Name()
		e.tmpFile.Close()
		os.Remove(name)
	}()

	wait := make(chan error)

	go func() {
		err := e.cmd.Wait()
		wait <- err
	}()

	select {
	case err := <-wait:
		return err
	case <-time.After(timeout):
		err := e.cmd.Process.Kill()
		return fmt.Errorf("killed process timeout, error of kill: %w", err)
	}
}

func (e Execution) Replace(str string) string {
	return strings.ReplaceAll(str, "GREU_TMP", e.tmpFile.Name())
}

func (e Execution) ReplaceMultiple(str ...string) []string {
	for i := range str {
		str[i] = e.Replace(str[i])
	}

	return str
}
