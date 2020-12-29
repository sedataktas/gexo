package main

import (
	"fmt"
	"golang.org/x/term"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	f := getFile()

	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Errorf("terminal raw error:%v", err)
	}
	defer term.Restore(fd, state)

	//term.NewTerminal(os.Stdin,">>>>>")

	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	os.Stdin.Write(f)
	var b = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		fmt.Print(string(b))
		if string(b) == "q" {
			exec.Command("stty", "-raw").Run()
		}
	}
}

func getFile() []byte {
	fileName := os.Args[1]

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	fByteArr, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return fByteArr
}
