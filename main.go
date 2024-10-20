package main

import (
	"fmt"
	// "os"
	"slices"
	"time"

	// "golang.org/x/term"
)

const gridWidth int8 = 10
const gridHeight int8 = 20

const blank string = "\033[41m "

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
		fmt.Println()
	}
}

// Function to update the grid in place
func render(grid [][]*string) {
	// ANSI escape sequence to move cursor to top-left corner
	fmt.Print("\033[H")
	printGrid(grid)
}

func main() {


	// reset view
	fmt.Print("\033[2J") // Clear screen first
	fmt.Print("\033[H")  // move to top left

	// printGrid(grid)

	oShape := shape{
		name:  "",
		block: "\033[46m ",
		grid:  [][]bool{{true, true}, {true, true}},
		top:   0,
		left:  4,
	}

	// Update grid values in a loop
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second) // Pause for a second
		grid := genGrid(oShape)

		render(grid)

		oShape.top += 1
	}
}

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

type point struct {
	x int
	y int
}

type shape struct {
	name  string
	block string
	grid  [][]bool

	top  int16
	left int16
}
