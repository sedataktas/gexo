package main

import (
	"fmt"
	"golang.org/x/term"
	"os"
)

func main() {
	fd := int(os.Stdin.Fd())

	state, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Errorf("terminal raw error:%v", err)
	}
	defer term.Restore(fd, state)

	newTerminal := term.NewTerminal(os.Stdin, ">")
	line, err := newTerminal.ReadLine()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(line)
}
