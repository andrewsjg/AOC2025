package main

import (
	"bufio"
	"fmt"
	"os"
)

func readInput(inputFile string) ([][]string, error) {
	output := [][]string{}

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return output, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()

		// convert line to a slice of characters
		characters := []string{}
		for _, char := range line {
			characters = append(characters, string(char))
		}
		output = append(output, characters)
		lineCount++

	}

	return output, nil

}

func solution(input [][]string) (part1 int, part2 int) {
	totalSplits := 0

	if len(input) == 0 {
		return 0, 0
	}

	width := len(input[0])
	pathTracker := make([]int, width)

	// Track current beam positions (true = beam present at this position)
	currentBeamPath := make([]bool, width)

	// Find starting position in the first row
	startIndex := findIndices(input[0], "S")
	if len(startIndex) > 0 {
		currentBeamPath[startIndex[0]] = true
		pathTracker[startIndex[0]] = 1
	}

	// Process each subsequent row
	for lineNo := 1; lineNo < len(input); lineNo++ {
		line := input[lineNo]
		nextBeamPath := make([]bool, width)

		// For each position, check if there's a beam
		for i := 0; i < width; i++ {
			if currentBeamPath[i] {
				if line[i] == "^" {
					// Beam hits a splitter
					totalSplits++

					pathTracker[i-1] += pathTracker[i]
					pathTracker[i+1] += pathTracker[i]
					pathTracker[i] = 0

					// Add left beam path
					if i > 0 {
						nextBeamPath[i-1] = true
					}
					// Add right beam path
					if i < width-1 {
						nextBeamPath[i+1] = true
					}

				} else {
					// No splitter - beam continues straight down
					nextBeamPath[i] = true
				}
			}
		}

		currentBeamPath = nextBeamPath
	}

	totalPaths := 0

	for _, paths := range pathTracker {
		totalPaths += paths
	}

	return totalSplits, totalPaths
}

func main() {
	inputFile := "input.txt"
	input, err := readInput(inputFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	totalSpits, totalPaths := solution(input)
	fmt.Printf("Part 1 Total Spits: %d\nPart 2 Total Paths: %d\n", totalSpits, totalPaths)

}

func findIndices(chars []string, target string) []int {
	indices := []int{}
	for i, char := range chars {
		if char == target {
			indices = append(indices, i)
		}
	}
	return indices
}
