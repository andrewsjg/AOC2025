package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IngredientRange struct {
	start int
	end   int
}

type IngredientRanges []IngredientRange
type IngredientList []int

func readInput(inputFile string) (IngredientRanges, IngredientList, error) {
	iRange := IngredientRanges{}
	iList := IngredientList{}

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return iRange, iList, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	rangeSection := true
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			rangeSection = false
			continue
		}

		if rangeSection {
			// Process ranges
			sRangeStart := strings.Split(line, "-")[0]
			sRangeEnd := strings.Split(line, "-")[1]

			rangeStart, err := strconv.Atoi(sRangeStart)

			if err != nil {
				return iRange, iList, err
			}

			rangeEnd, err := strconv.Atoi(sRangeEnd)

			if err != nil {
				return iRange, iList, err
			}

			iRange = append(iRange, IngredientRange{start: rangeStart, end: rangeEnd})

		} else {
			// Process ingredient list
			ingredient, err := strconv.Atoi(line)

			if err != nil {
				return iRange, iList, err
			}
			iList = append(iList, ingredient)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}

	// Sort iRange by start value
	for i := 0; i < len(iRange)-1; i++ {
		for j := i + 1; j < len(iRange); j++ {
			if iRange[i].start > iRange[j].start {
				// Swap
				iRange[i], iRange[j] = iRange[j], iRange[i]
			}
		}
	}

	//Sort iList
	for i := 0; i < len(iList)-1; i++ {
		for j := i + 1; j < len(iList); j++ {
			if iList[i] > iList[j] {
				// Swap
				iList[i], iList[j] = iList[j], iList[i]
			}
		}
	}

	return iRange, iList, nil
}

func part1(ranges IngredientRanges, ingredients IngredientList) (int, error) {

	freshIngredients := 0

	// First normalise the ranges
	newRanges := normaliseRanges(ranges)

	for _, ingredient := range ingredients {
		inRange := false
		for _, r := range newRanges {
			if ingredient >= r.start && ingredient <= r.end {
				inRange = true
				break
			}
		}

		if inRange {
			freshIngredients++
		}
	}

	return freshIngredients, nil
}

func part2(iRanges IngredientRanges) (int, error) {
	freshIngredients := 0

	// First normalise the ranges
	newRanges := normaliseRanges(iRanges)

	for _, freshrange := range newRanges {

		//fmt.Printf("Range start: %d Range end: %d\n", freshrange.start, freshrange.end)
		freshIngredients += (freshrange.end - freshrange.start + 1)
	}

	return freshIngredients, nil

}

func normaliseRanges(ranges IngredientRanges) IngredientRanges {

	/*
		3-5
		10-14
		16-20
		12-18

		Normalizes to:
		3-5
		10-20

	*/

	outputRanges := IngredientRanges{}
	outputRanges = append(outputRanges, ranges[0])

	for rangeNum, r := range ranges {
		// ignore first range as already added
		if rangeNum == 0 {
			continue
		}
		lastIdx := len(outputRanges) - 1

		if r.start <= outputRanges[lastIdx].end {
			outputRanges[lastIdx].end = max(outputRanges[lastIdx].end, r.end)
		} else {
			outputRanges = append(outputRanges, r)
		}

	}

	return outputRanges

}

func main() {
	iRange, iList, err := readInput("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	result, err := part1(iRange, iList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in part1: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 1 - Fresh ingredient count: %d\n", result)

	result, err = part2(iRange)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in part1: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 2 - Fresh ingredient count: %d\n", result)

}
