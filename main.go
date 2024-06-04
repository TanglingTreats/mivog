package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

type Pos struct {
	x int
	y int
}

type Cursor struct {
	postion Pos
}

// Escape sequences
const (
	ESCAPE  = "\x1B"
	SIGTERM = "\x03"
)

// Control sequences
const (
	TermErase = "\033[2J"
	CusrMvTop = "\033[%d;%dH" // template string
)

// Terminal states
const (
	NORMAL = iota
	INSERT
	VISUAL
	COMMAND
)

// Cursor styles
const (
	CusrBlk  = "\033[0 q"
	CusrLine = "\033[5 q"
)

// Normal Mode command
const (
	NNewline = 'o'
	NInsert  = 'i'
	NAppend  = 'a'
)

var termState = NORMAL

func main() {
	// Set terminal to raw
	state, err := term.GetState(0)
	if err != nil {
		log.Fatalln("failed to get state", err)
	}

	fmt.Println("make raw")
	_, err = term.MakeRaw(0)
	if err != nil {
		log.Fatalln("setting stdin to raw: ", err)
	}

	// Restore terminal
	defer func() {
		if err := term.Restore(0, state); err != nil {
			log.Println("warning, failed to restore terminal:", err)
		}
	}()

	fmt.Print(TermErase)
	moveCusrAbs(1, 2)

	reader := bufio.NewReader(os.Stdin)

	for {
		cursorState(termState)
		key, _, error := reader.ReadRune()
		if error != nil {
			fmt.Println("Error met")
			os.Exit(1)
		}

		if string(key) == SIGTERM {
			break
		}

		if string(key) == ESCAPE {
			termState = NORMAL
		}

		switch termState {
		case NORMAL:
			if key == NInsert {
				termState = INSERT
			}
		case INSERT:
			fmt.Print("\033[3 q")
			fmt.Print(string(key))
		case VISUAL:
		case COMMAND:
		}

	}
}

func moveCusrAbs(x int, y int) {
	fmt.Printf(CusrMvTop, x, y)
}

func moveCusrRel(x int, y int) {

}

func cursorState(state int) {
	switch state {
	case NORMAL:
		fmt.Print(CusrBlk)
	case INSERT:
		fmt.Print(CusrLine)
	case VISUAL:
	case COMMAND:
	}
}
