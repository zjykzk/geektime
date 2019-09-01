package geektimedl

import "fmt"

func ereaseCurrentLine() {
	fmt.Print("\x1b[2K\r")
}

func moveUp(i int) {
	fmt.Printf("\x1b[%dA", i)
}

func moveDown(i int) {
	fmt.Printf("\x1b[%dB", i)
}

func hideCursor() {
	fmt.Print("\x1b[?25l")
}

func showCursor() {
	fmt.Print("\x1b[?25h")
}
