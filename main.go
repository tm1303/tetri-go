package main

import (
	"fmt"
	// "image/color"
	"os"
	"slices"
	"sync"
	"time"

	"golang.org/x/term"
)

const (
	gridWidth  int = 10
	gridHeight int = 20
	gridBuffer int = 4

	showBuffer = true
)

var (
	buffer = black
	blank  = white
)

type point struct {
	x     int
	y     int
	color *string
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

func playPointsOk(playShape shape, binGrid []point, vMov int, hMov int, rotMov int) bool {

	renderPoints := []point{}

	rotIndex := playShape.gridIndex + rotMov
	if rotIndex < 0 {
		rotIndex = len(playShape.grids) - 1
	}
	if rotIndex >= len(playShape.grids) {
		rotIndex = 0
	}

	for shapeXIndex, row := range playShape.grids[rotIndex] {
		for shapeYIndex, v := range row {
			if v {
				renderPoints = append(renderPoints, point{
					x:     shapeXIndex + int(playShape.top+vMov),
					y:     shapeYIndex + int(playShape.left+hMov),
					color: &playShape.block,
				})
			}
		}
	}

	for _, rp := range renderPoints {

		if rp.y < 0 || rp.y > gridWidth-1 {
			return false
		}

		if rp.x > gridHeight+gridBuffer-1 {
			return false
		}

		if slices.ContainsFunc(binGrid, func(bp point) bool {
			return bp.x == rp.x && bp.y == rp.y
		}) {
			return false
		}
	}

	return true
}

func combinePoints(playShape shape, binGrid []point) []point {
	renderPoints := slices.Clone(binGrid)

	for shapeXIndex, row := range playShape.grids[playShape.gridIndex] {
		for shapeYIndex, v := range row {
			if v {
				renderPoints = append(renderPoints, point{
					x:     shapeXIndex + int(playShape.top),
					y:     shapeYIndex + int(playShape.left),
					color: &playShape.block,
				})
			}
		}
	}

	return renderPoints
}

// Generate grid based on the current shape position
func genGrid(playShape shape, binGrid []point) [][]*string {
	renderPoints := combinePoints(playShape, binGrid)

	grid := make([][]*string, gridHeight+gridBuffer)
	for rowIndex := range grid {
		grid[rowIndex] = make([]*string, gridWidth)
		for colIndex := range grid[rowIndex] {
			rpi := slices.IndexFunc(renderPoints, func(rp point) bool {
				return rowIndex == rp.x && colIndex == rp.y
			})
			if rpi > -1 {
				grid[rowIndex][colIndex] = renderPoints[rpi].color
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

	binGrid := make([]point, (gridHeight+gridBuffer)*gridWidth)

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
				if !playPointsOk(playShape, binGrid, 0, -1, 0) {
					continue
				}
				playShape.left--
				grid := genGrid(playShape, binGrid)
				render(grid)
			case 'd': // Move right
				if !playPointsOk(playShape, binGrid, 0, 1, 0) {
					continue
				}
				playShape.left++
				grid := genGrid(playShape, binGrid)
				render(grid)
			case 's': // Move down
				if !playPointsOk(playShape, binGrid, 1, 0, 0) {
					continue
				}
				playShape.top++
				grid := genGrid(playShape, binGrid)
				render(grid)
			case 'e': // rotate
				if !playPointsOk(playShape, binGrid, 0, 0, 1) {
					continue
				}
				playShape.gridIndex++
				if playShape.gridIndex >= len(playShape.grids) {
					playShape.gridIndex = 0
				}
				grid := genGrid(playShape, binGrid)
				render(grid)
			case 'q': // rotate
				if !playPointsOk(playShape, binGrid, 0, 0, -1) {
					continue
				}
				playShape.gridIndex--
				if playShape.gridIndex < 0 {
					playShape.gridIndex = len(playShape.grids) - 1
				}
				grid := genGrid(playShape, binGrid)
				render(grid)
			case 27: // Quit
				alive = false
				return
			}
		}
	}()

	// Update grid values in a loop
	for alive {

		grid := genGrid(playShape, binGrid)
		render(grid)
		time.Sleep(1 * time.Second) // Pause for a second

		if !playPointsOk(playShape, binGrid, 1, 0, 0) {

			binGrid = combinePoints(playShape, binGrid) // get in the bin

			shapeLibIndex++
			if shapeLibIndex >= len(shapeLib) {
				shapeLibIndex = 0
			}
			playShape = shapeLib[shapeLibIndex]
			playShape.top = 0
		}

		dropShape(&playShape) // Move down every loop iteration

	}

	wg.Wait() // Wait for the goroutine to finish
}

func dropShape(playShape *shape) {
	playShape.top++
}
