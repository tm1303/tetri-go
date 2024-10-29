package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

const (
	gridWidth  int = 10
	gridHeight int = 20
	gridBuffer int = 4
)

var (
	blank = white
)

type renderPoint struct {
	x     int
	y     int
	color *string
}

type playUnit struct {
	y     int
	color *string
}

type bin map[int][]playUnit

// Main function
func main() {

	wsLogger := initWsLog()
	log.Logger = log.Output(wsLogger)

	// Reset view
	fmt.Print("\033[2J") // Clear screen first
	fmt.Print("\033[H")  // Move to top left

	// Enable raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
			panic(err)
		}
	}()

	quitChan := make(chan interface{})
	inputChan := make(chan string)

	handleInput(quitChan, inputChan)
	gameLoop(quitChan, inputChan)
}
