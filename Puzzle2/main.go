package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func readInput(inputFile string) (string, error) {

	output := ""

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
		output = output + line
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}

	return output, nil
}

func part1(input string) (int, error) {

	result := 0
	invalidNums := []int{}
	// split the input into individual sequences around the commmas
	sequences := strings.Split(input, ",")

	for _, seq := range sequences {
		// split the sequeneces into start and end around the hyphen
		sL := strings.Split(seq, "-")[0]
		sU := strings.Split(seq, "-")[1]

		L, err := strconv.Atoi(sL)

		if err != nil {
			return result, err
		}
		U, err := strconv.Atoi(sU)

		if err != nil {
			return result, err
		}

		// Ignore sequences where both L and U have an odd number of digits
		if !((countDigits(L)%2 != 0) && (countDigits(U)%2 != 0)) {

			// increment L until it has an even number of digits
			for (countDigits(L)%2 != 0) && L <= U {
				L++
			}

			patternLen := countDigits(L) / 2
			multiplier := int(math.Pow(10, float64(patternLen))) + 1

			P1 := (L + multiplier - 1) / multiplier
			P2 := (U + multiplier - 1) / multiplier

			for P := P1; P <= P2; P++ {
				invalidNum := P * multiplier

				if (countDigits(invalidNum)%2 == 0) && (invalidNum >= L) && (invalidNum <= U) {
					invalidNums = append(invalidNums, invalidNum)
				}
			}

			//fmt.Printf("Invalid: %d, patternLen: %d multiplier: %d P1: %d, P2: %d\n\n", invalidNums, patternLen, multiplier, P1, P2)

		}
	}

	for _, num := range invalidNums {
		result += num
	}

	return result, nil
}

func part2(input string) (int, error) {

	result := 0 //invalidNums := []int{}
	invalidNums := map[int]bool{}

	// split the input into individual sequences around the commmas
	sequences := strings.Split(input, ",")

	for _, seq := range sequences {
		// split the sequeneces into start and end around the hyphen
		sL := strings.Split(seq, "-")[0]
		sU := strings.Split(seq, "-")[1]

		L, err := strconv.Atoi(sL)

		if err != nil {
			return result, err
		}
		U, err := strconv.Atoi(sU)

		if err != nil {
			return result, err
		}

		LD := countDigits(L)
		UD := countDigits(U)

		for d := LD; d <= UD; d++ {

			lower := max(L, (int(math.Pow(10, float64(d-1)))))
			upper := min(U, (int(math.Pow(10, float64(d))) - 1))

			for div := 1; div <= d/2; div++ {

				if d%div != 0 {
					continue
				}

				multiplier := (int(math.Pow(10, float64(d))) - 1) / (int(math.Pow(10, float64(div))) - 1)

				P1 := (lower + multiplier - 1) / multiplier
				P2 := (upper + multiplier - 1) / multiplier

				for P := P1; P <= P2; P++ {
					invalidNum := P * multiplier

					if (invalidNum >= L) && (invalidNum <= U) {

						invalidNums[invalidNum] = true
					}

				}
			}
		}

	}

	for num := range invalidNums {
		result += num
	}
	return result, nil
}

func countDigits(num int) int {
	if num == 0 {
		return 1
	}
	count := 0
	for num != 0 {
		num /= 10
		count++
	}
	return count
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

	fmt.Printf("Part 1 Result: %d\n", result)

	result, err = part2(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in part1: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Part 2 Result: %d\n", result)

}
