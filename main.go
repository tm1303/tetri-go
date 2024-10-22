package main

import (
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"golang.org/x/term"
)

const (
	gridWidth  int = 10
	gridHeight int = 20
	blank          = white
	red            = "\033[41m " //[]byte{keyEscape, '[', '3', '1', 'm'},
	green          = "\033[42m " //[]byte{keyEscape, '[', '3', '2', 'm'},
	yellow         = "\033[43m " //[]byte{keyEscape, '[', '3', '3', 'm'},
	blue           = "\033[44m " //[]byte{keyEscape, '[', '3', '4', 'm'},
	magenta        = "\033[45m " //[]byte{keyEscape, '[', '3', '5', 'm'},
	cyan           = "\033[46m " //[]byte{keyEscape, '[', '3', '6', 'm'},
	white          = "\033[47m " //[]byte{keyEscape, '[', '3', '7', 'm'},
)

type point struct {
	x int
	y int
}

type shapeGrid [][]bool

type shape struct {
	name      string
	block     string
	grids     []shapeGrid
	gridIndex int
	top       int
	left      int
}

var oShape = shape{
	name:      "",
	block:     red,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{true, true},
			{true, true},
		},
	},
	top:  0,
	left: 4,
}

var lShape = shape{
	name:      "",
	block:     green,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{false, true, false},
			{false, true, false},
			{false, true, true},
		},
		{
			{false, false, false},
			{true, true, true},
			{true, false, false},
		},
		{
			{true, true, false},
			{false, true, false},
			{false, true, false},
		},
		{
			{false, false, true},
			{true, true, true},
			{false, false, false},
		},
	},
	top:  0,
	left: 4,
}

var jShape = shape{
	name:      "",
	block:     yellow,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{false, true, false},
			{false, true, false},
			{true, true, false},
		},
		{
			{false, false, false},
			{true, true, true},
			{false, false, false},
		},
		{
			{false, true, true},
			{false, true, false},
			{false, true, false},
		},
		{
			{false, false, false},
			{true, true, true},
			{false, false, true},
		},
	},
	top:  0,
	left: 4,
}

var iShape = shape{
	name:      "",
	block:     blue,
	gridIndex: 0,
	grids: []shapeGrid{
		{
			{false, false, false, false},
			{true, true, true, true},
			{false, false, false, false},
			// {false, false, false, false},
		},
		{
			{false, true, false},
			{false, true, false},
			{false, true, false},
			{false, true, false},
		},
		// {
		// 	{false, false, false, false},
		// 	{false, false, false, false},
		// 	{true, true, true, true},
		// 	{false, false, false, false},
		// },
		// {
		// 	{false, true, false, false},
		// 	{false, true, false, false},
		// 	{false, true, false, false},
		// 	{false, true, false, false},
		// },
	},
	top:  0,
	left: 4,
}

var shapeLib = []shape{
	iShape,
	jShape,
	lShape,
	oShape,
}

// Function to print the grid
func printGrid(grid [][]*string) {
	for _, row := range grid {
		for _, val := range row {
			if val == nil {
				fmt.Printf("%s", blank)
			} else {
				fmt.Printf("%s", *val)
			}
		}
		fmt.Print("\r\n")
	}
}

// Function to update the grid in place
func render(grid [][]*string) {
	fmt.Print("\033[H") // ANSI escape to move cursor to top-left corner
	printGrid(grid)
}

// Generate grid based on the current shape position
func genGrid(playShape shape) [][]*string {
	renderPoints := []point{}

	for shapeXIndex, row := range playShape.grids[playShape.gridIndex] {
		for shapeYIndex, v := range row {
			if v {
				renderPoints = append(renderPoints, point{
					x: shapeXIndex + int(playShape.top),
					y: shapeYIndex + int(playShape.left),
				})
			}
		}
	}

	grid := make([][]*string, gridHeight)
	for rowIndex := range grid {
		grid[rowIndex] = make([]*string, gridWidth)
		for colIndex := range grid[rowIndex] {
			if slices.ContainsFunc(renderPoints, func(rp point) bool {
				return rowIndex == rp.x && colIndex == rp.y
			}) {
				grid[rowIndex][colIndex] = &playShape.block
				continue
			}
			grid[rowIndex][colIndex] = nil
		}
	}

	return grid
}

// Main function
func main() {
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

	var wg sync.WaitGroup
	wg.Add(1)

	shapeLibIndex := 0
	playShape := shapeLib[shapeLibIndex]

	alive := true
	// Goroutine to handle key presses
	go func() {
		defer wg.Done()
		for {
			var buf [1]byte
			_, err := os.Stdin.Read(buf[:]) // Read a single byte
			if err != nil {
				break // Exit on error
			}
			switch buf[0] {
			case 'a': // Move left
				if playShape.left > 0 {
					playShape.left--
					grid := genGrid(playShape)
					render(grid)
				}
			case 'd': // Move right
				if playShape.left < int(gridWidth-2) { // 2 for shape width
					playShape.left++
					grid := genGrid(playShape)
					render(grid)
				}
			case 's': // Move down
				if playShape.top < int(gridHeight-2) { // 2 for shape height
					playShape.top++
					grid := genGrid(playShape)
					render(grid)
				}
			case 'e': // rotate
				playShape.gridIndex++
				if playShape.gridIndex >= len(playShape.grids) {
					playShape.gridIndex = 0
				}
				grid := genGrid(playShape)
				render(grid)
			case 'q': // rotate
				playShape.gridIndex--
				if playShape.gridIndex < 0 {
					playShape.gridIndex = len(playShape.grids) - 1
				}
				grid := genGrid(playShape)
				render(grid)
			case 27: // Quit
				alive = false
				return
			}
		}
	}()

	// Update grid values in a loop
	for alive {

		grid := genGrid(playShape)
		render(grid)
		time.Sleep(1 * time.Second)      // Pause for a second
		playShape.top++                  // Move down every loop iteration
		if playShape.top >= gridHeight { // Reset shape position for demo purposes
			shapeLibIndex++
			if shapeLibIndex >= len(shapeLib) {
				shapeLibIndex = 0
			}
			playShape = shapeLib[shapeLibIndex]
			playShape.top = 0
		}
	}

	wg.Wait() // Wait for the goroutine to finish
}
