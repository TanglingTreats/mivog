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

// ASCII table
const (
	ETX = 0x03 //      End of Text
	LF  = 0x0A // '\n' Line Feed
	CR  = 0x0D // '\r' Carriage Return
	ESC = 0x1B // '\e' Escape
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

	NLeft  = "h"
	NUp    = "j"
	NDown  = "k"
	NRight = "l"
)

var termState = NORMAL

func main() {
	buf := make([]rune, 0, 2)

	// Set terminal to raw
	state, err := term.MakeRaw(0)
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

		if key == ETX {
			fmt.Print("\n\r")
			fmt.Println(string(buf))
			break
		}

		if key == ESC {
			termState = NORMAL
		}

		switch termState {
		case NORMAL:
			if key == NInsert {
				termState = INSERT
			}
		case INSERT:
			if key == CR {
				fmt.Print(string(LF) + string(CR))
				buf = append(buf, LF)
				buf = append(buf, CR)
			} else {
				buf = append(buf, key)
				fmt.Print(string(key))
			}
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
