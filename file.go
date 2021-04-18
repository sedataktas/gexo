package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
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
		return
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
