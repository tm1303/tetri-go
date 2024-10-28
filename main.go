package main

import (
	"fmt"
	// "log"
	// "image/color"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

const (
	gridWidth  int = 10
	gridHeight int = 20
	gridBuffer int = 4

	// showBuffer = true
)

var (
	// buffer = black
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

func playPointsOk(playShape shape, binGrid map[int][]playUnit, vMov int, hMov int, rotMov int) bool {

	rotIndex := playShape.gridIndex + rotMov
	if rotIndex < 0 {
		rotIndex = len(playShape.grids) - 1
	}
	if rotIndex >= len(playShape.grids) {
		rotIndex = 0
	}

	renderPoints := []renderPoint{}

	for shapeXIndex, row := range playShape.grids[rotIndex] {
		for shapeYIndex, v := range row {
			if v {
				renderPoints = append(renderPoints, renderPoint{
					x:     shapeXIndex + int(playShape.top+vMov),
					y:     shapeYIndex + int(playShape.left+hMov),
					color: &playShape.block,
				})
			}
		}
	}

	binPoints := make([]renderPoint, 0, len(binGrid))
	for x, playUnits := range binGrid {
		for _, playUnit := range playUnits {
			binPoints = append(binPoints, renderPoint{
				x:     x,
				y:     playUnit.y,
				color: playUnit.color,
			})
		}
	}

	for _, rp := range renderPoints {

		if rp.y < 0 || rp.y > gridWidth-1 {
			return false
		}

		if rp.x > gridHeight-1 {
			return false
		}

		if slices.ContainsFunc(binPoints, func(bp renderPoint) bool {
			return bp.x == rp.x && bp.y == rp.y
		}) {
			return false
		}
	}

	return true
}

func combinePoints(playShape shape, binGrid map[int][]playUnit) map[int][]playUnit {
	//renderPoints := maps.Clone(binGrid)
	renderPoints := make(map[int][]playUnit, 0)
	for i, binPoints := range binGrid {
		renderPoints[i] = slices.Clone(binPoints)
	}

	for shapeXIndex, row := range playShape.grids[playShape.gridIndex] {

		for shapeYIndex, v := range row {
			if v {
				newX := shapeXIndex + int(playShape.top)
				renderPoints[newX] = append(renderPoints[newX], playUnit{
					// x:     newX,
					y:     shapeYIndex + int(playShape.left),
					color: &playShape.block,
				})
			}
		}
	}

	return renderPoints
}

func tidyBin(binGrid map[int][]playUnit) (map[int][]playUnit, int) {

	newBin := map[int][]playUnit{}
	removedCount := 0

	keys := make([]int, 0, len(binGrid))
	for k := range binGrid {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	log.Debug().Msgf("keys %v", keys)

	for i := len(keys) - 1; i >= 0; i-- {
		k := keys[i]
		log.Debug().Msgf("key %d", k)
		binRow := binGrid[k]
		// for i, binRow := range binGrid {
		if len(binRow) >= gridWidth {
			log.Debug().Msgf("bin row %d", k)
			removedCount++
		} else {
			newBin[k+removedCount] = binRow
		}
	}

	if removedCount > 0 {
		log.Debug().Msgf("new bin  %v", newBin)
	}

	return newBin, removedCount
}

// Generate grid based on the current shape position
func genGrid(playShape shape, binGrid map[int][]playUnit) [][]*string {
	combined := combinePoints(playShape, binGrid)

	// renderPoints := make([]playUnit, 0, len(combined))
	// for _, v := range combined {
	// 	renderPoints = append(renderPoints, v...)
	// }

	grid := make([][]*string, gridHeight)
	for rowIndex := range grid {
		grid[rowIndex] = make([]*string, gridWidth)
		for colIndex := range grid[rowIndex] {
			rpi := slices.IndexFunc(combined[rowIndex], func(rp playUnit) bool {
				return colIndex == rp.y //rowIndex == rp.x && colIndex == rp.y
			})
			if rpi > -1 {
				grid[rowIndex][colIndex] = combined[rowIndex][rpi].color
				continue
			}
			grid[rowIndex][colIndex] = nil
		}
	}

	return grid
}

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

	var wg sync.WaitGroup
	wg.Add(1)

	shapeLibIndex := 0
	playShape := shapeLib[shapeLibIndex]
	playShape.top = -gridBuffer

	binGrid := make(map[int][]playUnit, 0)

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

	pauseFactor := time.Duration(1000)
	score := 0
	// Update grid values in a loop
	for alive {
		score++

		grid := genGrid(playShape, binGrid)
		render(grid)
		fmt.Printf("%sscore: %d\r\n", white, score)
		time.Sleep(pauseFactor * time.Millisecond) // Pause for a second

		if !playPointsOk(playShape, binGrid, 1, 0, 0) {
			score=score+2
			binGrid = combinePoints(playShape, binGrid) // get in the bin
			removedCount := 0
			binGrid, removedCount = tidyBin(binGrid)
			if removedCount > 0 {
				score = (4 * removedCount * removedCount) + score
				pauseFactor = time.Duration(0.95 * float64(pauseFactor))
			}

			shapeLibIndex++
			if shapeLibIndex >= len(shapeLib) {
				shapeLibIndex = 0
			}
			playShape = shapeLib[shapeLibIndex]
			playShape.top = -gridBuffer

			log.Debug().Msgf("new shape %s\n", playShape.name)

		}

		dropShape(&playShape) // Move down every loop iteration

	}

	wg.Wait() // Wait for the goroutine to finish
}

func dropShape(playShape *shape) {
	playShape.top++
}
