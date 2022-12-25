package main

import (
	"bytes"
	"fmt"
	"os"
)

type FuncHandler func(line []byte) (LineHandler, error)

func (p FuncHandler) HandleLine(line []byte) (LineHandler, error) {
	return p(line)
}

func (p FuncHandler) End() error {
	return nil
}

func passthroutHandler(line []byte) (LineHandler, error) {
	if bytes.HasPrefix(line, []byte(".GNUPLOT")) {
		fmt.Fprintln(os.Stderr, "  Changing to command")
		return NewCommandHandler(), nil
	}

	fmt.Print(string(line))
	return FuncHandler(passthroutHandler), nil
}
