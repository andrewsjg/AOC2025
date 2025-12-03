package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Read the input into a slice of strings so it can be used by both part1 and part2
func readInput(inputFile string) ([]string, error) {

	output := []string{}
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return output, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Process each line here
		output = append(output, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}

	return output, nil
}

func solution(input []string, part2 bool) {

	position := 50
	zeroCount := 0
	for _, line := range input {

		// process the component parts of each line
		direction := line[0:1]
		distance, err := strconv.Atoi(line[1:])

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting string to int: %v\n", err)
			os.Exit(1)
		}

		if direction == "R" {

			if part2 {
				zeroCount += (position + distance) / 100
			}

			position = (position + distance) % 100

		}

		if direction == "L" {

			if part2 {
				if position > 0 {
					zeroCount += (distance + 100 - position) / 100
				} else {
					zeroCount += distance / 100
				}
			}

			position = (position - distance) % 100

			if position < 0 {
				position += 100
			}

		}

		if !part2 {
			if position == 0 {
				zeroCount++
			}
		}
	}

	fmt.Printf("Final Position: %d, Zero Count: %d\n", position, zeroCount)

}

func main() {

	input, err := readInput("input.txt")
	if err != nil {
		os.Exit(1)
	}

	// Part 1
	solution(input, false)

	// Part 2
	solution(input, true)
}
