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
		E.row.size = len(t)
		E.row.bytes = &t
	}

	E.numRows = 1
}

func getFile(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	fByteArr, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return fByteArr
}
