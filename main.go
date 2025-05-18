package main

import (
	"fmt"
	"io"
	"math"
	"os"

	"golang.org/x/term"
)

var lastAnswer float64 = math.NaN()

func main() {
	fd := int(os.Stdin.Fd())
	origState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rw := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	terminal := term.NewTerminal(rw, ">")
	for {
		input, err := terminal.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		tokens := Tokenize(input)
		parser := NewParser(tokens)
		ast, err := parser.Expression()
		if err != nil {
			switch err.(type) {
			case ClearError:
				fmt.Fprint(terminal, "\x1b[2J\x1b[H")
				continue
			case ExitError:
				term.Restore(fd, origState)
				os.Exit(0)
			default:
				fmt.Fprintln(terminal, err.Error())
				continue
			}
		}
		lastAnswer, err = ast.Eval()
		if err != nil {
			fmt.Fprintln(terminal, err)
		} else {
			fmt.Fprintln(terminal, lastAnswer)
		}
	}
	term.Restore(fd, origState)
	fmt.Println()
}

