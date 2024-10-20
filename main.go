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
	gridWidth  int16 = 10
	gridHeight int16 = 20
	blank            = "\033[41m "
)

type point struct {
	x int
	y int
}

type shape struct {
	name  string
	block string
	grid  [][]bool
	top   int16
	left  int16
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

	for shapeXIndex, row := range playShape.grid {
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

	// Create a shape to be rendered
	oShape := shape{
		name:  "",
		block: "\033[46m ",
		grid:  [][]bool{{true, true}, {true, true}},
		top:   0,
		left:  4,
	}

	var wg sync.WaitGroup
	wg.Add(1)

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
				if oShape.left > 0 {
					oShape.left--
					grid := genGrid(oShape)
					render(grid)
				}
			case 'd': // Move right
				if oShape.left < int16(gridWidth-2) { // 2 for shape width
					oShape.left++
					grid := genGrid(oShape)
					render(grid)
				}
			case 's': // Move down
				if oShape.top < int16(gridHeight-2) { // 2 for shape height
					oShape.top++
					grid := genGrid(oShape)
					render(grid)
				}
			case 'q': // Quit
				alive = false
				return
			}
		}
	}()

	// Update grid values in a loop
	for alive {
		grid := genGrid(oShape)
		render(grid)
		time.Sleep(1 * time.Second)   // Pause for a second
		oShape.top++                  // Move down every loop iteration
		if oShape.top >= gridHeight { // Reset shape position for demo purposes
			oShape.top = 0
		}
	}

	wg.Wait() // Wait for the goroutine to finish
}
