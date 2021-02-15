package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	ArrowLeft int = iota + 1000
	ArrowRight
	ArrowUp
	ArrowDown
)

const version = "0.0.1"

var buf string

func editorReadKey(reader *bufio.Reader) int {
	// read one byte
	c, err := reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			fmt.Println("END OF FILE")
		}
	}

	var seq [2]byte
	// Arrow keys begin with escape character
	// then comes '[', then followed 'A', 'B', 'C', 'D'
	if c == '\x1b' {
		// read 2 more bytes
		_, err := reader.Read(seq[:])
		if err != nil {
			panic(err)
		}

		if seq[0] == '[' {
			switch seq[1] {
			case 'A':
				return ArrowUp
			case 'B':
				return ArrowDown
			case 'C':
				return ArrowRight
			case 'D':
				return ArrowLeft
			}
		}
		return '\x1b'
	}
	return int(c)
}

func editorProcessKeypress(c int) {
	switch c {
	case ctrlKey('q'):
		clearEntireScreen()
		getCursorToBegin()
		disableRawMode()
		os.Exit(1)
	case ArrowUp, ArrowDown, ArrowLeft, ArrowRight:
		editorMoveCursor(c)
		break
	}
}

func editorRefreshScreen() {
	hideCursor()
	//clearEntireScreen()
	getCursorToBegin()

	editorDrawRows()
	setCursorPosition()

	showCursor()

	_, err := fmt.Print(buf)
	if err != nil {
		panic(err)
	}
}

func hideCursor() {
	// l --> reset mode
	buf += "\x1b[?25l"
}

func showCursor() {
	// h --> set mode
	buf += "\x1b[?25h"
}

func clearEntireScreen() {
	// The first byte is \x1b, which is the escape character, or 27 in decimal
	// The other three bytes are [2J.
	// We are writing an escape sequence to the terminal.
	//Escape sequences always start with an escape character (27) followed by a [ character.
	//Escape sequences instruct the terminal to do various text formatting tasks,
	//such as coloring text, moving the cursor around, and clearing parts of the screen.

	//We are using the J command (Erase In Display) to clear the screen.
	//Escape sequence commands take arguments, which come before the command.
	//In this case the argument is 2, which says to clear the entire screen.
	//<esc>[1J would clear the screen up to where the cursor is,
	//and <esc>[0J would clear the screen from the cursor up to the end of the screen.
	//Also, 0 is the default argument for J, s
	//o just <esc>[J by itself would also clear the screen from the cursor to the end.
	buf += "\x1b[2J"
}

func getCursorToBegin() {
	//  <esc>[2J command left the cursor at the bottom of the screen.
	// We have to reposition it at the top-left corner so
	//that we’re ready to draw the editor interface from top to bottom.
	// For this we use H command for take the cursor to the first row and first column
	//_, err := writer.Write([]byte("\x1b[H"))
	buf += "\x1b[H"
}

func editorDrawRows() {
	for i := 0; i < E.screenRows; i++ {
		if i == E.screenRows/3 {
			welcomeMsg := fmt.Sprintf("gexo editor -- version %s", version)
			buf += "~" + " "
			padding := (E.screenCols - len(welcomeMsg)) / 2
			for j := 0; j <= padding; j++ {
				buf += " "
			}
			buf += welcomeMsg
		} else {
			// write tilde sign to each row
			buf += "~"
		}

		// erase in line : https://vt100.net/docs/vt100-ug/chapter3.html#EL, default : 0
		buf += "\x1b[K"
		if i < E.screenRows-1 {
			buf += "\r\n"
		}
	}
}

func setCursorPosition() {
	buf += fmt.Sprintf("\x1b[%d;%dH", E.cy+1, E.cx+1)
}

func editorMoveCursor(key int) {
	switch key {
	case ArrowLeft:
		if E.cx != 0 {
			E.cx--
		}
		break
	case ArrowRight:
		if E.cx != E.screenCols-1 {
			E.cx++
		}
		break
	case ArrowUp:
		if E.cy != 0 {
			E.cy--
		}
		break
	case ArrowDown:
		if E.cy != E.screenRows-1 {
			E.cy++
		}
		break
	}
}
