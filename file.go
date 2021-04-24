package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

func editorOpen(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		if t != "" {
			byteArray := []byte(t)
			editorInsertRow(E.numRows, byteArray)
		}
	}

	E.dirty = 0
}

func fileSave() {
	if E.fileName == "" {
		E.fileName = string(editorPrompt("Save as: %s (ESC to cancel)"))
		if E.fileName == "" {
			editorSetStatusMessage("Save aborted")
			return
		}
	}

	str := editorRowsToString()

	err := ioutil.WriteFile(E.fileName, []byte(str), 0644)
	if err != nil {
		editorSetStatusMessage(fmt.Sprintf("Can't save! I/O error: %v", err))
		panic(err)
	}

	E.dirty = 0
	editorSetStatusMessage(fmt.Sprintf("%d bytes written to disk", len(str)))
}

func editorPrompt(prompt string) []byte {
	bufSize := 128
	var buf []byte

	for {
		editorSetStatusMessage(fmt.Sprintf(prompt, buf))
		editorRefreshScreen()

		c := editorReadKey(bufio.NewReader(os.Stdin))
		if c == DeleteKEy ||
			c == ctrlKey('h') ||
			c == BackSpace {
			if len(buf) != 0 {
				buf = buf[:len(buf)-1]
			}
		} else if c == '\x1b' {
			editorSetStatusMessage("")
			buf = nil
			return nil
		} else if c == '\r' {
			if len(buf) != 0 {
				editorSetStatusMessage("")
				return buf
			}
		} else if !unicode.IsControl(rune(c)) && c < 128 {
			if len(buf) == bufSize-1 {
				bufSize *= 2
			}
			buf = append(buf, byte(c))
		}
	}
}

func editorFind() {
	query := editorPrompt("Search: %s (ESC to cancel)")
	if query == nil {
		return
	}

	// TODO : buraya esc tuşuna basılmadığı müddetce ve enter tuşuna basıldıktan
	// sonra da devam edecek bir mantık yazılmalı
	// yukarı aşağı tuşlarıyla bir sonraki aramanın yapılabilmesi için

	for i := 0; i < E.numRows; i++ {
		matchedIndex := strings.Index(string(E.row[i].bytes), string(query))
		if matchedIndex != -1 {
			E.cy = i
			E.cx = matchedIndex
			E.rowOff = E.numRows

			for _, _ = range query {
				E.row[i].highlights = insert(E.row[i].highlights, matchedIndex, byte(HlMatch))
				matchedIndex++
			}
			break
		}
	}
}
