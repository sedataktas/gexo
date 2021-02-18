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

		r := Erow{
			size:  len(t),
			bytes: &t,
		}

		E.row = append(E.row, r)
		E.numRows++
	}
}
