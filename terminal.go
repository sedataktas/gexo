package main

import (
	"golang.org/x/sys/unix"
	"os"
)

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
	// &^ = bit clear (AND NOT)
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

	// On some systems, when you type Ctrl-V,
	// the terminal waits for you to type another character and then sends that character literally
	// IEXTEN diasble this feature
	raw.Lflag &^= unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN

	// By default, Ctrl-S and Ctrl-Q are used for software flow control.
	// Ctrl-S stops data from being transmitted to the terminal until you press Ctrl-Q.
	// IXON disable ctrl+s and ctrl+q

	// Ctrl-M is weird: it’s being read as 10, when we expect it to be read as 13,
	// since it is the 13th letter of the alphabet, and Ctrl-J already produces a 10. What else produces 10? The Enter key does.
	// It turns out that the terminal is helpfully translating any carriage returns (13, '\r')
	// inputted by the user into newlines (10, '\n')

	// ISTRIP causes the 8th bit of each input byte to be stripped, meaning it will set it to 0.

	// INPCK enables parity checking

	// When BRKINT is turned on, a break condition will cause a SIGINT signal
	// to be sent to the program, like pressing Ctrl-C.
	raw.Iflag &^= unix.ICRNL | unix.IXON | unix.BRKINT | unix.INPCK | unix.ISTRIP

	// CS8 is not a flag, it is a bit mask with multiple bits,
	//which we set using the bitwise-OR (|) operator unlike all the flags we are turning off.
	//It sets the character size (CS) to 8 bits per byte.
	raw.Cflag &^= unix.CS8

	// It turns out that the terminal does a similar translation on the output side.
	// It translates each newline ("\n") we print into a carriage return followed by a newline ("\r\n").
	// The terminal requires both of these characters in order to start a new line of text.
	// The carriage return moves the cursor back to the beginning of the current line,
	// and the newline moves the cursor down a line, scrolling the screen if necessary
	// OPOST diasble this
	raw.Oflag &^= unix.OPOST

	// TODO : search in google
	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0
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
