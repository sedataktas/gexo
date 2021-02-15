package main

import (
	"bufio"
	"golang.org/x/sys/unix"
	"os"
)

type EditorConfig struct {
	screenRows  int
	screenCols  int
	origTermios unix.Termios
}

var (
	E EditorConfig
)

func init() {
	w, err := unix.IoctlGetWinsize(int(os.Stdin.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		panic(err)
	}

	E.screenCols = int(w.Col)
	E.screenRows = int(w.Row)
}

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
