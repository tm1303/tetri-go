package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
)

func dropShape(playShape *shape) {
	playShape.top++
}

func move(h int, v int, playShape *shape, binGrid *bin) {
	if !playPointsOk(playShape, binGrid, v, h, 0) {
		return
	}
	playShape.left=playShape.left+h
	playShape.top=playShape.top+v
	grid := updateGrid(playShape, binGrid)
	render(grid)
}

func rotate(r int, playShape *shape, binGrid *bin) {
	if !playPointsOk(playShape, binGrid, 0, 0, r) {
		return
	}
	playShape.gridIndex++
	if playShape.gridIndex >= len(playShape.grids) {
		playShape.gridIndex = 0
	}
	if playShape.gridIndex < 0 {
		playShape.gridIndex = len(playShape.grids) - 1
	}
	grid := updateGrid(playShape, binGrid)
	render(grid)
}

// func xgameLoop()
// {
	
// }

func gameLoop(quitChan chan interface{}, inputChan chan string) {

	pauseFactor := time.Duration(1000)
	score := 0

	shapeLibIndex := 0
	playShape := &shapeLib[shapeLibIndex]
	playShape.top = -gridBuffer

	binGrid := make(bin, 0)

	go func() {
		for {
			select {
			case action := <-inputChan:
				switch action {
				case "down":
					move(0, 1, playShape, &binGrid)
				case "left":
					move(-1, 0, playShape, &binGrid)
				case "right":
					move(1, 0, playShape, &binGrid)
				case "rotc":
					rotate(1, playShape, &binGrid)
				case "rota":
					rotate(-1, playShape, &binGrid)
				}
			}
		}
	}()

	for{

		select {
		case <-quitChan:
			log.Info().Msg("quit message recieved, killing loop")
			return
		case  <-time.After(pauseFactor * time.Millisecond):
			score++

			grid := updateGrid(playShape, &binGrid)
			render(grid)
			fmt.Printf("%sscore: %d\r\n", white, score)
			time.Sleep(pauseFactor * time.Millisecond) // Pause for a second

			if !playPointsOk(playShape, &binGrid, 1, 0, 0) {
				score = score + 2
				binGrid = combinePoints(playShape, &binGrid) // get in the bin
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
				playShape = &shapeLib[shapeLibIndex]
				playShape.top = -gridBuffer

				log.Debug().Msgf("new shape %s\n", playShape.name)

			}

			dropShape(playShape) // Move down every loop iteration
	}}
}


func playPointsOk(playShape *shape, binGrid *bin, vMov int, hMov int, rotMov int) bool {

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

	binPoints := make([]renderPoint, 0, len(*binGrid))
	for x, playUnits := range *binGrid {
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

func combinePoints(playShape *shape, binGrid *bin) bin {
	//renderPoints := maps.Clone(binGrid)
	renderPoints := make(bin, 0)
	for i, binPoints := range *binGrid {
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

func tidyBin(binGrid bin) (bin, int) {

	newBin := bin{}
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
func updateGrid(playShape *shape, binGrid *bin) [][]*string {
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
