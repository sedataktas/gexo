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
			byteArray = append(byteArray, '\000')
			r := Erow{
				size:  len(byteArray),
				bytes: byteArray,
			}

			E.row = append(E.row, r)
			E.numRows++
		}
	}
}
