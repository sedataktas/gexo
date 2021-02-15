package main

import (
	"bufio"
	"golang.org/x/sys/unix"
	"os"
)

type EditorConfig struct {
	cx          int
	cy          int
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

	// set cursor initial positions to 0
	E.cx = 0
	E.cy = 0

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

func ctrlKey(c byte) int {
	// The CTRL_KEY macro bitwise-ANDs a character with the value 00011111, in binary.
	ctrlKey := (c) & 0x1f
	return int(ctrlKey)
}
