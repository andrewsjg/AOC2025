package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readInput(inputFile string) ([][]int, error) {
	output := [][]int{}

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
		nums := []int{}

		for i := 0; i < len(line); i++ {

			num, err := strconv.Atoi(string(line[i]))
			if err != nil {
				return output, err
			}
			nums = append(nums, num)
		}

		output = append(output, nums)

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}

	return output, nil
}

func part1(input [][]int) (int, error) {
	result := 0

	for _, bank := range input {
		// fmt.Printf("Battery Bank: %v\n", bank)
		bankMax := 0
		highest := 0
		for battIndex, battJoltage := range bank {
			if battIndex < len(bank)-1 {
				if battJoltage > highest {
					highest = battJoltage
				}

				jolts := []int{highest, bank[battIndex+1]}
				if makeJolts(jolts) > bankMax {
					bankMax = makeJolts(jolts)
				}
			}

		}

		result = result + bankMax

	}

	return result, nil
}

// Part 1 approach couldn't be extended to part 2. Rethink the naive implementation

func part2(input [][]int) (int, error) {
	result := 0

	for _, bank := range input {

		enabled := []int{}
		available := len(bank) - 12

		for _, battJoltage := range bank {

			for available > 0 && len(enabled) > 0 && (enabled[len(enabled)-1] < battJoltage) {
				enabled = enabled[:len(enabled)-1]
				available--
			}

			enabled = append(enabled, battJoltage)
		}

		if len(enabled) >= 12 {
			enabled = enabled[:12]
		}

		/*
			fmt.Printf("Bank: %v\n", bank)
			fmt.Printf("Enabled: %v\n", enabled)
			fmt.Printf("jolts: %v\n\n", makeJolts(enabled))
		*/

		result = result + makeJolts(enabled)
	}

	return result, nil
}

func makeJolts(digits []int) int {
	result := 0
	for _, d := range digits {
		result = result*10 + d
	}
	return result
}

func main() {

	input, err := readInput("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	result, err := part1(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in part1: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1 - Total output joltage: %d\n\n", result)

	result, err = part2(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in part2: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 2 - Total output joltage: %d\n", result)

}
