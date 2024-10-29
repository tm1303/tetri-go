package main

import "fmt"

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