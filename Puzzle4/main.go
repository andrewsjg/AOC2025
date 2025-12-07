package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	row, col int
}

type Grid map[Point]bool

// Read the input, make it useful for both parts
func readInput(inputFile string) (Grid, error) {

	output := make(map[Point]bool)
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return output, err
	}
	defer file.Close()

	row := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row++

		for col, char := range line {
			if char == '@' {
				output[Point{row: row, col: col + 1}] = true
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		}
	}

	return output, nil
}

// Solve part 1
func part1(grid Grid) int {

	result := 0
	for gridLoc := range grid {
		if countNeighbours(grid, gridLoc) < 4 {
			result++
		}
	}

	return result
}

// Solve part 2 - Do it recursively
func part2(grid Grid, total int) int {

	removable := []Point{}
	for gridLoc := range grid {
		neighbourCount := countNeighbours(grid, gridLoc)

		if neighbourCount < 4 {
			removable = append(removable, gridLoc)
		}

	}

	// exit the recursion
	if len(removable) == 0 {
		return total
	}

	for _, loc := range removable {
		delete(grid, loc)
		total++
	}

	return part2(grid, total)
}

// The trick is to only add locations for the rolls (@).
// Learnt this in previous years.
func countNeighbours(grid Grid, gridLoc Point) int {
	neighbourCount := 0

	for gridRow := -1; gridRow <= 1; gridRow++ {
		for gridCol := -1; gridCol <= 1; gridCol++ {
			if gridRow == 0 && gridCol == 0 {
				continue
			}

			if grid[Point{gridLoc.row + gridRow, gridLoc.col + gridCol}] {
				neighbourCount++
			}
		}
	}

	return neighbourCount
}

func main() {
	grid, err := readInput("input.txt")
	if err != nil {
		os.Exit(1)
	}

	result := part1(grid)
	fmt.Printf("Part 1 - Total accessible rolls:\t%d\n", result)

	result = part2(grid, 0)
	fmt.Printf("Part 2 - Total removable rolls:\t\t%d\n", result)
}
