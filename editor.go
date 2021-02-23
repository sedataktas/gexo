package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	// set random number, but number is bigger than character's limit
	ArrowLeft int = iota + 1000
	ArrowRight
	ArrowUp
	ArrowDown
	PageUp
	PageDown
	HomeKey
	EndKey
	DeleteKEy
)

const version = "0.0.1"

var buf []byte

func editorReadKey(reader *bufio.Reader) int {
	// read one byte
	c, err := reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			fmt.Println("END OF FILE")
		}
	}

	var seq [3]byte
	// Arrow keys begin with escape character
	// then comes '[', then followed 'A', 'B', 'C', 'D'
	if c == '\x1b' {
		// read 2 more bytes
		_, err := reader.Read(seq[:])
		if err != nil {
			panic(err)
		}

		if seq[0] == '[' {
			// Page Up is sent as <esc>[5~ and Page Down is sent as <esc>[6~.
			// The Home key could be sent as <esc>[1~, <esc>[7~, <esc>[H, or <esc>OH.
			// Similarly, the End key could be sent as <esc>[4~, <esc>[8~, <esc>[F, or <esc>OF.
			// Delete key : It simply sends the escape sequence <esc>[3~
			if seq[1] >= '0' && seq[1] <= '9' {
				if &seq[2] == nil {
					return '\x1b'
				}

				if seq[2] == '~' {
					switch seq[1] {
					case '1':
						return HomeKey
					case '3':
						return DeleteKEy
					case '4':
						return EndKey
					case '5':
						return PageUp
					case '6':
						return PageDown
					case '7':
						return HomeKey
					case '8':
						return EndKey
					}
				}
			} else if seq[0] == 'O' {
				switch seq[1] {
				case 'H':
					return HomeKey
				case 'F':
					return EndKey
				}
			} else {
				switch seq[1] {
				case 'A':
					return ArrowUp
				case 'B':
					return ArrowDown
				case 'C':
					return ArrowRight
				case 'D':
					return ArrowLeft
				case 'H':
					return HomeKey
				case 'F':
					return EndKey
				}
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
	case PageDown:
		for i := 0; i < E.screenRows; i++ {
			editorMoveCursor(ArrowDown)
		}
	case PageUp:
		for i := 0; i < E.screenRows; i++ {
			editorMoveCursor(ArrowUp)
		}
	case HomeKey:
		E.cx = 0
		break
	case EndKey:
		E.cx = E.screenCols - 1
		break
	}
}

func editorRefreshScreen() {
	editorScroll()
	hideCursor()
	getCursorToBegin()

	editorDrawRows()
	setCursorPosition()

	showCursor()

	_, err := fmt.Print(string(buf))
	if err != nil {
		panic(err)
	}
}

func hideCursor() {
	// l --> reset mode
	buf = append(buf, []byte("\x1b[?25l")...)
}

func showCursor() {
	// h --> set mode
	buf = append(buf, []byte("\x1b[?25h")...)
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
	buf = append(buf, []byte("\x1b[2J")...)
}

func getCursorToBegin() {
	//  <esc>[2J command left the cursor at the bottom of the screen.
	// We have to reposition it at the top-left corner so
	//that we’re ready to draw the editor interface from top to bottom.
	// For this we use H command for take the cursor to the first row and first column
	//_, err := writer.Write([]byte("\x1b[H"))
	buf = append(buf, []byte("\x1b[H")...)
}

func editorDrawRows() {
	for i := 0; i < E.screenRows; i++ {
		fileRow := i + E.rowOff
		if fileRow >= E.numRows {
			if E.numRows == 0 && i == E.screenRows/3 {
				welcomeMsg := fmt.Sprintf("gexo editor -- version %s", version)
				buf = append(buf, '~', ' ')
				padding := (E.screenCols - len(welcomeMsg)) / 2
				for j := 0; j <= padding; j++ {
					buf = append(buf, ' ')
				}
				msgByteArr := []byte(welcomeMsg)
				buf = append(buf, msgByteArr...)
			} else {
				// write tilde sign to each row
				buf = append(buf, '~')
			}
		} else {
			len := E.row[fileRow].rsize - E.colOff
			if len < 0 {
				len = 0
			}
			if len > E.screenCols {
				len = E.screenCols
			}

			buf = append(buf, E.row[fileRow].render...)
		}

		// erase in line : https://vt100.net/docs/vt100-ug/chapter3.html#EL, default : 0

		buf = append(buf, []byte("\x1b[K")...)
		if i < E.screenRows-1 {
			buf = append(buf, []byte("\r\n")...)
		}
	}
}

func setCursorPosition() {
	cur := fmt.Sprintf("\x1b[%d;%dH", (E.cy-E.rowOff)+1, (E.rx-E.colOff)+1)
	buf = append(buf, []byte(cur)...)

}

func editorMoveCursor(key int) {
	var row *Erow
	if E.cy >= E.numRows {
		row = nil
	} else {
		row = &E.row[E.cy]
	}

	switch key {
	case ArrowLeft:
		if E.cx != 0 {
			E.cx--
		} else if E.cy > 0 {
			E.cy--
			E.cx = E.row[E.cy].size
		}
		break
	case ArrowRight:
		if row != nil && E.cx < row.size {
			E.cx++
		} else if row != nil && E.cx == row.size {
			E.cy++
			E.cx = 0
		}
		break
	case ArrowUp:
		if E.cy != 0 {
			E.cy--
		}
		break
	case ArrowDown:
		if E.cy < E.numRows {
			E.cy++
		}
		break
	}

	if E.cy >= E.numRows {
		row = nil
	} else {
		row = &E.row[E.cy]
	}

	rowLen := 0
	if row != nil {
		rowLen = row.size
	}

	if E.cx > rowLen {
		E.cx = rowLen
	}
}

func editorScroll() {
	E.rx = E.cx

	if E.cy < E.numRows {
		E.rx = editorRowCxToRx(&E.row[E.cy], E.cx)
	}
	if E.cy < E.rowOff {
		E.rowOff = E.cy
	}
	if E.cy >= E.rowOff+E.screenRows {
		E.rowOff = E.cy - E.screenRows + 1
	}

	if E.rx < E.colOff {
		E.colOff = E.rx
	}
	if E.rx >= E.colOff+E.screenCols {
		E.colOff = E.rx - E.screenCols + 1
	}
}

func editorAppendRow(bytes []byte, len int) {
	at := E.numRows

	E.row[at].size = len
	E.row[at].bytes = bytes
	E.row[at].bytes[len] = '\000'

	E.row[at].rsize = 0
	E.row[at].render = nil
	editorUpdateRow(&E.row[at])
}

func editorUpdateRow(row *Erow) {
	whiteSpace := 0
	for j := 0; j < row.size; j++ {
		if string(row.bytes[j]) == " " {
			if len(row.bytes) > j+3 {
				for i := 0; i < 3; i++ {
					if string(row.bytes[j]) == " " {
						j++
						whiteSpace++
					}
				}
			}
			if whiteSpace == 3 {
				row.render = append(row.render, ' ')
			}
			whiteSpace = 0
		} else {
			row.render = append(row.render, row.bytes[j])
		}
	}

	/*
		row.render = nil

		idx := 0
		for j :=0 ; j<row.size; j++ {
			if string(row.bytes[j]) == " " {
				if  len(row.bytes) > j+3{
					for i:=0; i<3 ;i++ {
						if string(row.bytes[j]) == " " {
							whiteSpace++
						}
					}
				}
				j=j+3
				if whiteSpace == 3 {
					row.render = append(row.render, ' ')
					for idx % 8 != 0 {
						row.render[idx+1] = ' '
					}
				}

			}else {
				if len(row.render) > idx+1 {
					row.render[idx+1] = row.bytes[j]
				} else {
					row.render = append(row.render, row.bytes[j])
				}

			}
		}
	*/
	//row.render[idx] = '\000'
	//row.rsize = idx
}

func editorRowCxToRx(row *Erow, cx int) int {
	rx := 0

	for j := 0; j < cx; j++ {
		if row.bytes[j] == '\t' {
			rx += 7 - (rx % 8)
		}
		rx++
	}
	return rx
}
