package main

import (
	"bufio"
	"golang.org/x/sys/unix"
	"os"
)

const ()

var origTermios unix.Termios

func main() {
	reader := bufio.NewReader(os.Stdin)

	defer clearEntireScreen()
	defer getCursorToBegin()
	defer disableRawMode()

	enableRawMode()
	for {
		editorRefreshScreen()
		c := editorReadKey(reader)
		editorProcessKeypress(c)
	}
}

func ctrlKey(c byte) byte {
	// The CTRL_KEY macro bitwise-ANDs a character with the value 00011111, in binary.
	return (c) & 0x1f
}
