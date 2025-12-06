package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func init() {
	registerDay("06a", day06a)
	registerDay("06b", day06b)
}

func day06a(args []string, rdr io.Reader) (string, error) {
	debug := false
	sum := 0
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSpace(string(in)), "\n")
	if len(lines) < 3 {
		return "", fmt.Errorf("unexpected number of input lines %d", len(lines))
	}
	opRow := []rune{}
	for opStr := range strings.FieldsSeq(lines[len(lines)-1]) {
		switch opStr {
		case "*":
			opRow = append(opRow, '*')
		case "+":
			opRow = append(opRow, '+')
		default:
			return "", fmt.Errorf("unknown operator %s in position %d", opStr, len(opRow)+1)
		}
	}
	numRows := [][]int{}
	for row := range len(lines) - 1 {
		numRows = append(numRows, []int{})
		for numStr := range strings.FieldsSeq(lines[row]) {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return "", fmt.Errorf("failed to parse number %s on row %d pos %d: %w", numStr, row+1, len(numRows[row])+1, err)
			}
			numRows[row] = append(numRows[row], num)
		}
		if len(numRows[row]) != len(opRow) {
			return "", fmt.Errorf("unexpected row length, row %d has %d entries, op row has %d entries", row, len(numRows[row]), len(opRow))
		}
	}
	for col, op := range opRow {
		add := 0
		switch op {
		case '*':
			add = 1
			for row := range len(numRows) {
				add *= numRows[row][col]
			}
		case '+':
			for row := range len(numRows) {
				add += numRows[row][col]
			}
		}
		if debug {
			fmt.Fprintf(os.Stderr, "from col %d adding %d\n", col, add)
		}
		sum += add
	}
	return fmt.Sprintf("%d", sum), nil
}

func day06b(args []string, rdr io.Reader) (string, error) {
	debug := false
	sum := 0
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(in), "\n")
	opRow := len(lines) - 1
	for strings.TrimSpace(lines[opRow]) == "" {
		opRow--
		if opRow < 2 {
			return "", fmt.Errorf("unexpected number of input lines, opRow = %d", opRow)
		}
	}
	opCur := byte('x')
	add := 0
	numStr := ""
	for col := range len(lines[opRow]) {
		switch lines[opRow][col] {
		case ' ':
			// continuation of existing opCur
		case '+':
			opCur = lines[opRow][col]
			add = 0
		case '*':
			opCur = lines[opRow][col]
			add = 1
		default:
			return "", fmt.Errorf("unexpected operation %c in column %d", lines[opRow][col], col)
		}
		// append each non-space entry of the column to the number
		for row := range opRow {
			if lines[row][col] != ' ' {
				numStr += string(lines[row][col])
			}
		}
		if numStr == "" {
			// end of the previous op
			if debug {
				fmt.Fprintf(os.Stderr, "from col ending %d adding %d\n", col-1, add)
			}
			sum += add
			add = 0
		} else {
			// convert column to number
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return "", fmt.Errorf("failed to convert column %d with value %s to number: %w", col, numStr, err)
			}
			// run the op on the add accumulator
			switch opCur {
			case '+':
				add += num
			case '*':
				add *= num
			default:
				return "", fmt.Errorf("unexpected operation %c when processing column %d", opCur, col)
			}
			numStr = ""
		}
	}
	// handle the last col
	if debug {
		fmt.Fprintf(os.Stderr, "from last col adding %d\n", add)
	}
	sum += add
	return fmt.Sprintf("%d", sum), nil
}
