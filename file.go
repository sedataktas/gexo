package main

import (
	"bufio"
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
			r := Erow{
				size:   len(byteArray),
				rsize:  0,
				bytes:  byteArray,
				render: nil,
			}

			E.row = append(E.row, r)
			E.numRows++
			editorUpdateRow(&E.row[E.numRows])
		}
	}
}
