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

const (
	ESCAPE  = "\x1B"
	SIGTERM = "\x03"
)

const (
	NORMAL = iota
	INSERT
	VISUAL
	COMMAND
)

const (
	CUSR_BLK  = "\033[0 q"
	CUSR_LINE = "\033[5 q"
)

const (
	N_NEWLINE = 'o'
	N_INSERT  = 'i'
	N_APPEND  = 'a'
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

	fmt.Print("\033[2J")
	fmt.Print("\033[1;2H")

	// buf := make([]byte, 1)
	reader := bufio.NewReader(os.Stdin)
	// scanner := bufio.NewScanner(os.Stdin)

	// go getInput(scanner, buf)

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
			if key == N_INSERT {
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

func cursorState(state int) {
	switch state {
	case NORMAL:
		fmt.Print(CUSR_BLK)
	case INSERT:
		fmt.Print(CUSR_LINE)
	case VISUAL:
	case COMMAND:
	}
}

func getInput(scanner *bufio.Scanner, buf []byte) {
	for scanner.Scan() {
		buf = scanner.Bytes()
		fmt.Print("\033[1F")

		fmt.Println(string(buf))

		if string(buf) == "quit" {
			os.Exit(0)
		}
	}
}
