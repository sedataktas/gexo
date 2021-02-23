package main

import (
	"bufio"
	"golang.org/x/sys/unix"
	"os"
)

// Erow stands for “editor row”,
//and stores a line of text as a pointer to character data and a length
type Erow struct {
	size   int
	rsize  int
	bytes  []byte
	render []byte
}

type EditorConfig struct {
	cx          int
	cy          int
	rx          int
	rowOff      int
	colOff      int
	screenRows  int
	screenCols  int
	numRows     int
	row         []Erow
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
	E.rx = 0
	E.rowOff = 0
	E.colOff = 0
	E.numRows = -1
	E.row = nil
	E.screenCols = int(w.Col)
	E.screenRows = int(w.Row)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	defer clearEntireScreen()
	defer getCursorToBegin()
	defer disableRawMode()

	enableRawMode()

	if len(os.Args) > 1 {
		fileName := os.Args[1]
		editorOpen(fileName)
	}

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
