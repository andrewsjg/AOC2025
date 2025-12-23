package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	aocutilites "AOC2025/aocutilities"
)

func readInput(inputFile string) ([][]int, []string, error) {
	file, err := os.Open(inputFile)

	operandsOutput := [][]int{}
	operationsOutput := []string{}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return operandsOutput, operationsOutput, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		tmpArr := aocutilites.SplitWords(line)
		numArr := []int{}

		for _, word := range tmpArr {

			if word != "+" && word != "-" && word != "*" && word != "/" {
				num, err := strconv.Atoi(word)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error converting string to int: %v\n", err)
					return operandsOutput, operationsOutput, err
				}
				numArr = append(numArr, num)
			} else {
				operationsOutput = append(operationsOutput, word)
			}
		}
		if len(numArr) > 0 {
			operandsOutput = append(operandsOutput, numArr)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}

	return operandsOutput, operationsOutput, err
}

func readInputPart2(inputFile string) ([][]string, error) {
	output := [][]string{}
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return output, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// read the line into an array of characters including spaces
		charArr := []string{}
		for _, char := range line {
			charArr = append(charArr, string(char))
		}
		output = append(output, charArr)

	}

	return output, nil
}

func part1(operands [][]int, operations []string) int {

	total := 0
	stackTotal := 0

	for currentCol, op := range operations {

		calcStack := aocutilites.Stack[int]{}

		for _, row := range operands {
			calcStack.Push(row[currentCol])
		}

		stackTotal = stackMaths(op, calcStack)

		total += stackTotal
		stackTotal = 0

	}

	return total
}

// Ok this code is a mess! Cobbled together iteratively to get to the answer.
// But it does the job.
func part2(inputFile string) int {
	total := 0

	output, err := readInputPart2(inputFile)
	if err != nil {
		return total
	}

	ops := output[len(output)-1]
	operations := []string{}

	for _, op := range ops {
		if op != " " {
			operations = append(operations, op)
		}
	}

	output = output[:len(output)-1]
	rows := len(output) - 1
	cols := len(output[0])

	// Temp Output to visualize
	/*
		for _, row := range output {
			fmt.Printf("%v\n", row)
		}
		fmt.Printf("%v\n\n", operations)
	*/
	opIndex := 0
	op := ""

	prevWasBlank := false
	calcStack := aocutilites.Stack[int]{}

	for col := 0; col < cols; col++ {
		isBlankCol := true
		for rowCount := 0; rowCount <= rows; rowCount++ {
			if output[rowCount][col] != " " {
				isBlankCol = false
				break
			}
		}

		if !isBlankCol {
			nums := []string{}
			for rowCount := 0; rowCount <= rows; rowCount++ {
				currentDigit := output[rowCount][col]
				if currentDigit != " " {
					nums = append(nums, currentDigit)
				}
			}

			if len(nums) > 0 {

				if opIndex < len(operations) {
					op = operations[opIndex]
				}

				tmpNums := []int{}

				for _, n := range nums {
					num, err := strconv.Atoi(n)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error converting string to int: %v\n", err)
						return total
					}
					tmpNums = append(tmpNums, num)
				}

				numval := combineDigits(tmpNums)
				calcStack.Push(numval)

			}

			prevWasBlank = false
		} else {
			// Hit a blank column - do the operation on the stack
			// and switch to next operation

			stackTotal := stackMaths(op, calcStack)

			total += stackTotal
			calcStack = aocutilites.Stack[int]{}

			if !prevWasBlank {
				opIndex++
			}
			prevWasBlank = true
		}
	}
	// Do the final sum
	stackTotal := stackMaths(op, calcStack)
	total += stackTotal

	return total
}

func main() {
	input := "input.txt"
	operands, operations, err := readInput(input)

	if err != nil {
		os.Exit(1)
	}

	total := part1(operands, operations)
	fmt.Printf("Grand Total Part 1: %d\n", total)

	// Took a different approach for part 2
	total = part2(input)
	fmt.Printf("Grand Total Part 2: %d\n", total)
}

// Should refactor this to enable the operations to be pushed onto the stack
// with the operands.
func stackMaths(op string, stack aocutilites.Stack[int]) int {
	total := 0

	for !stack.IsEmpty() {
		val, _ := stack.Pop()

		switch op {
		case "+":
			total += val
		case "-":
			total -= val
		case "*":
			if total == 0 {
				total = 1
			}
			total *= val
		case "/":
			if val != 0 {
				total /= val
			}
		}
	}

	return total
}

func combineDigits(digits []int) int {
	result := 0
	for _, digit := range digits {
		result = result*10 + digit
	}
	return result
}
