package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
	"unicode"
)

const (
	BackSpace = 127
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

const (
	HlNormal int = iota
	HlNumber
	HlMatch
)

const (
	version      = "0.0.1"
	minQuitTimes = 3
)

var buf []byte
var quitTimes = minQuitTimes

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
	case '\r':
		editorInsertNewline()
	case ctrlKey('q'):
		if E.dirty != 0 && quitTimes > 0 {
			editorSetStatusMessage(fmt.Sprintf("WARNING!!! File has unsaved changes. "+
				"Press Ctrl-Q %d more times to quit.", quitTimes))
			quitTimes--
			return
		}
		clearEntireScreen()
		getCursorToBegin()
		disableRawMode()
		os.Exit(1)
	case ctrlKey('s'):
		fileSave()
	case ArrowUp, ArrowDown, ArrowLeft, ArrowRight:
		editorMoveCursor(c)
		/*
			case PageDown:
				for i := 0; i < E.screenRows; i++ {
					editorMoveCursor(ArrowDown)
				}
			case PageUp:
				for i := 0; i < E.screenRows; i++ {
					editorMoveCursor(ArrowUp)
				}

		*/
	case PageDown, PageUp:
		if c == PageUp {
			E.cy = E.rowOff
		} else if c == PageDown {
			E.cy = E.rowOff + E.screenRows - 1
			if E.cy > E.numRows {
				E.cy = E.numRows
			}
		}

		for i := E.screenRows; i > 0; i-- {
			if c == PageUp {
				editorMoveCursor(ArrowUp)
			} else {
				editorMoveCursor(ArrowDown)
			}
		}
	case HomeKey:
		E.cx = 0
	case EndKey:
		E.cx = E.screenCols - 1
	case BackSpace, ctrlKey('h'), DeleteKEy:
		if c == DeleteKEy {
			editorMoveCursor(ArrowRight)
		}
		editorDelChar()
	case ctrlKey('f'):
		editorFind()
	case ctrlKey('l'):
	case '\x1b':
	default:
		editorInsertChar(c)
	}
	quitTimes = minQuitTimes
}

func editorRefreshScreen() {
	editorScroll()

	buf = nil
	hideCursor()
	getCursorToBegin()

	editorDrawRows()
	editorDrawStatusBar()
	editorDrawMessageBar()
	setCursorPosition()

	showCursor()

	_, err := fmt.Print(string(buf))
	if err != nil {
		panic(err)
	}
	buf = nil
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
	for i := 0; i < E.screenRows-1; i++ {
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
			//buf = append(buf, E.row[fileRow].bytes...)
			len := E.row[fileRow].size - E.colOff
			if len < 0 {
				len = 0
			}

			if len > E.screenCols {
				len = E.screenCols
			}

			currentColor := -1
			for i := 0; i < len; i++ {
				if E.row[fileRow].highlights[i] == byte(HlNormal) {
					if currentColor != -1 {
						// 39 is reset color
						// more information : https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
						buf = append(buf, []byte("\x1b[39m")...)
						currentColor = -1
					}
					buf = append(buf, E.row[fileRow].bytes[i])
				} else {
					color := syntaxToColor(int(E.row[fileRow].highlights[i]))
					if color != currentColor {
						colorFmt := fmt.Sprintf("\x1b[%dm", color)
						buf = append(buf, []byte(colorFmt)...)

						// 31 is red
						//buf = append(buf, []byte("\x1b[31m")...)
					}
					buf = append(buf, E.row[fileRow].bytes[i])
				}
				buf = append(buf, []byte("\x1b[39m")...)
			}
		}

		// erase in line : https://vt100.net/docs/vt100-ug/chapter3.html#EL, default : 0
		buf = append(buf, []byte("\x1b[K")...)
		buf = append(buf, []byte("\r\n")...)
	}
}

func setCursorPosition() {
	cur := fmt.Sprintf("\x1b[%d;%dH", (E.cy-E.rowOff)+1, (E.cx-E.colOff)+1)
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
	case ArrowRight:
		if row != nil && E.cx < row.size {
			E.cx++
		} else if row != nil && E.cx == row.size {
			E.cy++
			E.cx = 0
		}
	case ArrowUp:
		if E.cy != 0 {
			E.cy--
		}
	case ArrowDown:
		if E.cy < E.numRows {
			E.cy++
		}
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
	if E.cy < E.rowOff {
		E.rowOff = E.cy
	}
	if E.cy >= E.rowOff+E.screenRows {
		E.rowOff = E.cy - E.screenRows + 1
	}

	if E.cx < E.colOff {
		E.colOff = E.cx
	}
	if E.cx >= E.colOff+E.screenCols {
		E.colOff = E.cx - E.screenCols + 1
	}
}

func editorDrawStatusBar() {
	//The m command (Select Graphic Rendition) causes the text printed
	//after it to be printed with various possible attributes including
	//bold (1),
	//underscore (4),
	//blink (5), and
	//inverted colors (7). F
	//or example, you could specify all of these attributes
	//using the command <esc>[1;4;5;7m. An argument of 0 clears all attributes,
	//and is the default argument, so we use <esc>[m to go back to normal text formatting.
	buf = append(buf, "\x1b[7m"...)

	statusBarMsgLeft := getStatusBarMsgLeft()
	statusBarMsgRight := fmt.Sprintf("%d/%d", E.cy+1, E.numRows)
	rlen := len(statusBarMsgRight)

	len := len(statusBarMsgLeft)
	if len > E.screenCols {
		len = E.screenCols
	}
	for i := 0; i < len; i++ {
		buf = append(buf, statusBarMsgLeft[i])
	}

	for len < E.screenCols {
		if E.screenCols-len == rlen {
			buf = append(buf, statusBarMsgRight...)
			break
		} else {
			// draw a blank white status bar of inverted space characters
			buf = append(buf, " "...)
			len++
		}
	}

	//  <esc>[m switches back to normal formatting
	buf = append(buf, "\x1b[m"...)
	buf = append(buf, "\r\n"...)
}

func editorDrawMessageBar() {
	buf = append(buf, "\x1b[K"...)
	msgLen := len(E.statusMsg)
	if msgLen > E.screenCols {
		msgLen = E.screenCols
	}

	t := time.Since(E.statusMsgTime).Seconds()
	if msgLen != 0 && t < 5 {
		buf = append(buf, E.statusMsg...)
	}
}

func editorSetStatusMessage(str ...string) {
	var buffer bytes.Buffer

	for _, s := range str {
		buffer.WriteString(s)
	}

	E.statusMsg = buffer.String()
	E.statusMsgTime = time.Now()
}

func editorInsertChar(c int) {
	if E.cy == E.numRows {
		editorInsertRow(E.numRows, []byte(""))
	}

	editorRowInsertChar(&E.row[E.cy], E.cx, c)
	//editorUpdateSyntax(&E.row[E.cy])
	E.cx++
}

func editorDelChar() {
	if E.cy == E.numRows {
		return
	}

	if E.cx == 0 && E.cy == 0 {
		return
	}
	if E.cx > 0 {
		editorRowDelChar(&E.row[E.cy], E.cx-1)
		E.cx--
	} else {
		E.cx = E.row[E.cy-1].size
		editorRowAppendString(&E.row[E.cy-1], E.row[E.cy].bytes)
		editorDelRow(E.cy)
		E.cy--
	}
}

func editorDelRow(at int) {
	if at < 0 || at >= E.numRows {
		return
	}

	E.row = append(E.row[:at], E.row[at+1:]...)

	E.numRows--
	E.dirty++
}

func editorRowAppendString(row *Erow, byteArray []byte) {
	row.bytes = append(row.bytes, byteArray...)
	row.size += len(byteArray)
	E.dirty++
	editorUpdateSyntax(row)
}

func editorRowDelChar(row *Erow, at int) {
	if at < 0 || at >= row.size {
		return
	}

	row.size--
	row.bytes = remove(row.bytes, at)
	E.dirty++
	editorUpdateSyntax(row)
}

func editorRowInsertChar(row *Erow, at, c int) {
	if at < 0 || at > row.size {
		at = row.size
	}

	row.size++
	row.bytes = insert(row.bytes, at, byte(c))
	row.highlights = insertHighlightSyntax(row.highlights, at, byte(c))
	E.dirty++
	//row.bytes = append(row.bytes, byte(c))
}

func editorInsertRow(at int, byteArray []byte) {
	if at < 0 || at > E.numRows {
		return
	}

	r := Erow{
		size:       len(byteArray),
		bytes:      byteArray,
		highlights: nil,
	}
	editorUpdateSyntax(&r)

	if len(E.row)-1 <= at {
		E.row = append(E.row, r)
	} else {
		E.row = insertRow(E.row, at, r)
	}

	E.numRows++
	E.dirty++
}

func editorInsertNewline() {
	if E.cx == 0 {
		editorInsertRow(E.cy, []byte(""))
	} else {
		row := &E.row[E.cy]
		editorInsertRow(E.cy+1, row.bytes[E.cx:])
		row = &E.row[E.cy]
		row.size = E.cx
		E.row = append(E.row, *row)
	}

	E.cy++
	E.cx = 0
}

func editorRowsToString() string {
	var buf string
	for _, row := range E.row {
		for _, r := range row.bytes {
			buf += string(r)
		}
		buf += string('\n')
	}
	return buf
}

func editorUpdateSyntax(row *Erow) {
	for _, b := range row.bytes {
		if unicode.IsDigit(rune(b)) {
			row.highlights = append(row.highlights, byte(HlNumber))
		} else {
			row.highlights = append(row.highlights, byte(HlNormal))
		}
	}
}

func syntaxToColor(hl int) int {
	switch hl {
	case HlNumber:
		return 31 // foreground red
	case HlMatch:
		return 34
	default:
		return 37 // foreground white
	}
}

func insert(a []byte, index int, value byte) []byte {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func insertHighlightSyntax(a []byte, index int, value byte) []byte {
	color := -1
	if unicode.IsDigit(rune(value)) {
		color = HlNumber
	} else {
		color = HlNormal
	}

	if len(a) == index { // nil or empty slice or after last element
		return append(a, byte(color))
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = byte(color)
	return a
}

func insertRow(rows []Erow, index int, row Erow) []Erow {
	if len(rows) == index { // nil or empty slice or after last element
		return append(rows, row)
	}
	rows = append(rows[:index+1], rows[index:]...) // index < len(a)
	rows[index] = row
	return rows
}

func remove(slice []byte, s int) []byte {
	return append(slice[:s], slice[s+1:]...)
}

func getStatusBarMsgLeft() string {
	statusBarMsgLeft := ""
	if E.fileName == "" {
		if E.dirty == 0 {
			statusBarMsgLeft = fmt.Sprintf("%.20s - %d lines %s", "[No Name]", E.numRows, "")

		} else {
			statusBarMsgLeft = fmt.Sprintf("%.20s - %d lines %s", "[No Name]", E.numRows, "(modified)")
		}
	} else {
		if E.dirty == 0 {
			statusBarMsgLeft = fmt.Sprintf("%.20s - %d lines %s", E.fileName, E.numRows, "")
		} else {
			statusBarMsgLeft = fmt.Sprintf("%.20s - %d lines %s", E.fileName, E.numRows, "(modified)")
		}
	}
	return statusBarMsgLeft
}
