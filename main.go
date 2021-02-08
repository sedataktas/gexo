package main

import (
	"bufio"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
)

const ()

var origTermios unix.Termios

func main() {
	defer disableRawMode()
	enableRawMode()
	reader := bufio.NewReader(os.Stdin)
	for {
		c := editorReadKey(reader)
		editorProcessKeypress(c)
	}
}

func ctrlKey(c byte) byte {
	// The CTRL_KEY macro bitwise-ANDs a character with the value 00011111, in binary.
	return (c) & 0x1f
}

func editorReadKey(reader *bufio.Reader) byte {
	// read one byte
	c, err := reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			fmt.Println("END OF FILE")
		}
	}
	return c
}

func editorProcessKeypress(c byte) {
	switch c {
	case ctrlKey('q'):
		disableRawMode()
		os.Exit(1)
	}

	/*
		if unicode.IsControl(rune(c)) {
			fmt.Printf("%d\r\n", c)
		} else {
			fmt.Printf("%d ('%c')\r\n", c, c)
		}
	*/
}

func editorRefreshScreen() {

}
