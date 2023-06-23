package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
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

	for {
		line, readErr := r.ReadBytes('\n')

		// Detect start of command
		if cmd, detected := cfg.DetectOpenCommand(line); inputState == Pass && detected {
			args := make([]string, len(cmd.Args))
			copy(args, cmd.Args)

			execution, err = NewExecution(cmd.Cmd, args...)
			if err != nil {
				panic(err)
			}

			currentCommand = cmd
			inputState = Command

			// Write prefix
			fmt.Fprintln(execution, execution.Replace(currentCommand.InputPrefix))
			fmt.Println(execution.Replace(currentCommand.ReplaceOpenTag))

			continue
		}

		// Detect end of command
		if close := currentCommand.DetectCloseCommand(line); inputState == Command && close {
			// Write postfix
			fmt.Fprintln(execution, execution.Replace(currentCommand.InputPostfix))

			err := execution.Exit(time.Second)
			if err != nil {
				panic(err)
			}

			// Write close tag after flush
			fmt.Println(execution.Replace(currentCommand.ReplaceCloseTag))

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

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			panic(err)
		}

	}
}
