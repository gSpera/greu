package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type LineHandler interface {
	HandleLine(line []byte) (LineHandler, error)
	End() error
}

type InputState int

const (
	Pass InputState = iota
	Command
)

func main() {
	cfgFile := flag.String("cfg", "greu.yml", "config file")

	r := bufio.NewReader(os.Stdin)
	inputState := Pass
	currentCommand := new(CommandDefinition)
	var execution Execution

	flag.Parse()

	cfg, err := LoadConfig(*cfgFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot decode config:", err)
		os.Exit(1)
	}
	fmt.Println(cfg)

	for {
		line, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		// Detect start of command
		if cmd, detected := cfg.DetectOpenCommand(line); inputState == Pass && detected {
			args := strings.Split(cmd.Cmd, " ")[1:]
			execCmd := strings.Split(cmd.Cmd, " ")[0]
			embeddedArgs := strings.Split(string(line), " ")[1:]
			args = append(args, embeddedArgs...)
			execution, err = NewExecution(execCmd, args...)
			if err != nil {
				panic(err)
			}

			currentCommand = cmd
			inputState = Command

			// Write prefix
			fmt.Fprintln(execution, currentCommand.InputPrefix)
			fmt.Println(currentCommand.ReplaceOpenTag)

			continue
		}

		// Detect end of command
		if close := currentCommand.DetectCloseCommand(line); inputState == Command && close {
			// Write postfix
			fmt.Fprintln(execution, currentCommand.InputPostfix)

			err := execution.Exit(time.Second)
			if err != nil {
				panic(err)
			}

			// Write close tag after flush
			fmt.Println(currentCommand.ReplaceCloseTag)

			inputState = Pass
			currentCommand = nil
			continue
		}

		// Redirect the input to correct writer
		switch inputState {
		case Pass:
			_, err = os.Stdout.Write(line)
		case Command:
			_, err = execution.Write(line)
			if err != nil {
				panic(err)
			}
		}

		if err != nil {
			panic(err)
		}
	}
}
