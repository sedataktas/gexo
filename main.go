package main

import (
	"bufio"
	"fmt"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"unicode"

	"io"

	"os"
)

var origTermios unix.Termios

func main() {
	enableRawMode()
	reader := bufio.NewReader(os.Stdin)
	for {
		var err error
		// read one byte
		c, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				fmt.Println("END OF FILE")
			}
		}
		// press q to quit.
		if c == 'q' {
			os.Exit(0)
		}

		if unicode.IsControl(rune(c)) {
			fmt.Printf("%d\n", c)
		} else {
			fmt.Printf("%d ('%c')\n", c, c)
		}

	}
}

// Raw mode vs canonical mode
// https://unix.stackexchange.com/questions/21752/what-s-the-difference-between-a-raw-and-a-cooked-device-driver
func enableRawMode() {
	// The termios functions describe a general terminal interface that
	// is provided to control asynchronous communications ports.
	// what s termios more info : https://blog.nelhage.com/2009/12/a-brief-introduction-to-termios/
	origTermios, err := unix.IoctlGetTermios(int(os.Stdin.Fd()), unix.TIOCGETA)
	if err != nil {
		panic(err)
	}

	raw := *origTermios
	// disable echoing
	// Lflag is a local flag
	// more info about lflag : https://blog.nelhage.com/2009/12/a-brief-introduction-to-termios-termios3-and-stty/
	// ECHO is a bitflag, defined as 00000000000000000000000000001000 in binary.
	// We use the bitwise-NOT operator (~) on this value to get 11111111111111111111111111110111.
	// We then bitwise-AND this value with the flags field, which forces the fourth bit in the flags field to become 0,
	// and causes every other bit to retain its current value

	// There is an ICANON flag that allows us to turn off canonical mode.
	// This means we will finally be reading input byte-by-byte, instead of line-by-line.

	// By default, Ctrl-C sends a SIGINT signal to the current process which causes it to terminate,
	// and Ctrl-Z sends a SIGTSTP signal to the current process which causes it to suspend.
	// ISIG disables ctrl+c, ctrl+z and ctrl+y and read as ASCII bytes.
	raw.Lflag &^= unix.ECHO | unix.ICANON | unix.ISIG

	// Apply terminal attributes
	err = unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TIOCSETA, &raw)
	if err != nil {
		panic(err)
	}
}

func disableRawMode() {
	err := unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TIOCSETA, &origTermios)
	if err != nil {
		panic(err)
	}
}

func getFile() []byte {
	fileName := os.Args[1]

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	fByteArr, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return fByteArr
}
