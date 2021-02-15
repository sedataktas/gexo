package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const version = "0.0.1"

var buf string

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
		clearEntireScreen()
		getCursorToBegin()
		disableRawMode()
		os.Exit(1)
		break
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
	hideCursor()
	//clearEntireScreen()
	getCursorToBegin()

	editorDrawRows()

	getCursorToBegin()
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
