package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type CommandHandler struct {
	cmd       *exec.Cmd
	pipeRead  io.ReadCloser
	pipeWrite io.WriteCloser
}

func NewCommandHandler() *CommandHandler {
	command := &CommandHandler{}
	cmd := exec.Command("gnuplot")
	command.cmd = cmd
	pr, pw := io.Pipe()
	command.pipeRead = pr
	command.pipeWrite = pw
	cmd.Stdin = pr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	return command
}

func (c *CommandHandler) HandleLine(line []byte) (LineHandler, error) {
	if bytes.HasPrefix(line, []byte(".GNUPLOT")) {
		fmt.Fprintln(os.Stderr, "Closing command")
		err := c.End()
		return FuncHandler(passthroutHandler), err
	}

	fmt.Fprintln(os.Stderr, "Command:", string(line))
	_, err := c.pipeWrite.Write(line)
	return c, err
}

func (c *CommandHandler) End() error {
	c.pipeWrite.Close()
	c.cmd.Wait()
	return nil
}
