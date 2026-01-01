package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

func init() {
	registerDay("10a", day10a)
	registerDay("10b", day10b)
}

type day10Machine struct {
	indicatorLights  []bool
	buttonSchematics [][]int
	joltageReq       []int
}

func day10a(args []string, rdr io.Reader) (string, error) {
	debug := true
	sum := 0
	machines, err := day10ReadMachines(rdr)
	if err != nil {
		return "", err
	}

	for i, machine := range machines {
		best := math.MaxInt
		// each button should be pressed at most one time, try each combination (2^buttons)
		// pressing twice would undo the press
		for toggleBits := range 1 << len(machine.buttonSchematics) {
			// track the result of this combination, and how many buttons needed to be pressed
			result := make([]bool, len(machine.indicatorLights))
			toggles := 0
			buttonCur := 0
			for toggleBits != 0 {
				toggleVal := toggleBits % 2
				if toggleVal == 1 {
					// button was pressed, invert the specified indicator lights
					toggles++
					for _, indicator := range machine.buttonSchematics[buttonCur] {
						result[indicator] = !result[indicator]
					}
				}
				// shift the remaining bits, and track which button we are looking at
				toggleBits >>= 1
				buttonCur++
			}
			if toggles < best && day10IndicatorsMatch(result, machine.indicatorLights) {
				best = toggles
			}
		}
		if best == math.MaxInt {
			return "", fmt.Errorf("failed to find a toggle sequence for machine %v", machine)
		}
		if debug {
			fmt.Fprintf(os.Stderr, "machine %d requires %d toggles\n", i, best)
		}
		sum += best
	}

	return fmt.Sprintf("%d", sum), nil
}

func day10b(args []string, rdr io.Reader) (string, error) {
	debug := true
	sum := 0
	machines, err := day10ReadMachines(rdr)
	if err != nil {
		return "", err
	}

	for i, machine := range machines {
		matrix := make([][]int, len(machine.joltageReq))
		for row, joltage := range machine.joltageReq {
			matrix[row] = make([]int, len(machine.buttonSchematics)+1)
			matrix[row][len(machine.buttonSchematics)] = joltage
		}
		for col, buttons := range machine.buttonSchematics {
			for _, row := range buttons {
				matrix[row][col] = 1
			}
		}
		day10GaussianElimination(matrix, 0, 0)
		if debug {
			fmt.Fprintf(os.Stderr, "machine %d matrix:\n", i)
			for _, row := range matrix {
				fmt.Fprintf(os.Stderr, "  %v\n", row)
			}
		}

		// process the matrix to find the smallest sum of button values where each value is a integer >= 0
		// fail if any value is fractional or negative
		buttons := make([]struct {
			bound, dependent bool
			cur, max, row    int
		}, len(machine.buttonSchematics))
		for col := range len(buttons) {
			rowMatch := -1
			allZero := true
			for row := range len(matrix) {
				if matrix[row][col] != 0 {
					for otherCol := range len(matrix[row]) - 1 {
						if otherCol != col && matrix[row][otherCol] != 0 {
							allZero = false
							break
						}
					}
					if rowMatch != -1 {
						rowMatch = -2 // variable value
						break
					}
					rowMatch = row
				}
			}
			switch rowMatch {
			case -1:
				// no matches
				buttons[col].bound = true // bound to 0
			case -2:
				// multiple matches, set the max, the min is always 0
				limit := math.MaxInt
				for _, joltage := range machine.buttonSchematics[col] {
					limit = min(limit, machine.joltageReq[joltage])
				}
				buttons[col].max = limit
			default:
				if allZero {
					// exact match
					buttons[col].cur = matrix[rowMatch][len(matrix[rowMatch])-1]
					if matrix[rowMatch][col] != 1 {
						if buttons[col].cur%matrix[rowMatch][col] != 0 {
							// fail on a fractional value
							return "", fmt.Errorf("failed to compute exact value for %d using row %v", col, matrix[rowMatch])
						}
						buttons[col].cur /= matrix[rowMatch][col]
					}
					buttons[col].bound = true
				} else {
					buttons[col].row = rowMatch
					buttons[col].dependent = !allZero
				}
			}
		}
		iterate := func() bool {
			for col := range len(buttons) {
				if buttons[col].bound {
					continue
				}
				buttons[col].cur++
				if buttons[col].cur <= buttons[col].max {
					return true
				}
				buttons[col].cur = 0
			}
			return false
		}
		best := math.MaxInt
		ok := true
		joltages := make([]int, len(machine.joltageReq))
		for ; ok; ok = iterate() {
			valid := true
			// calculate dependent buttons
			for col := len(buttons) - 1; col >= 0; col-- {
				if buttons[col].dependent {
					val := matrix[buttons[col].row][len(matrix[buttons[col].row])-1]
					// deduct other solved buttons from the value
					for otherCol, matrixVal := range matrix[buttons[col].row][:len(matrix[buttons[col].row])-1] {
						if otherCol != col && matrixVal != 0 {
							val -= matrixVal * buttons[otherCol].cur
						}
					}
					if val < 0 || val%matrix[buttons[col].row][col] != 0 {
						valid = false
						break
					} else {
						buttons[col].cur = val / matrix[buttons[col].row][col]
					}
				}
			}
			if !valid {
				continue
			}
			// sum joltages to check for a valid solution
			for j := range joltages {
				joltages[j] = 0
			}
			for b, buttonSpec := range machine.buttonSchematics {
				for _, j := range buttonSpec {
					joltages[j] += buttons[b].cur
				}
			}
			for j := range joltages {
				if joltages[j] != machine.joltageReq[j] {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}
			// if valid, update the best count
			count := 0
			for _, button := range buttons {
				count += button.cur
			}
			if count <= best {
				best = count
				if debug {
					fmt.Fprintf(os.Stderr, "machine %d, found solution of %d using [", i, count)
					for b, button := range buttons {
						if b != 0 {
							fmt.Fprint(os.Stderr, ", ")
						}
						fmt.Fprintf(os.Stderr, "%d", button.cur)
					}
					fmt.Fprint(os.Stderr, "]\n")
				}
			}
		}
		if best == math.MaxInt {
			return "", fmt.Errorf("failed to find a solution")
		}
		sum += best
	}
	return fmt.Sprintf("%d", sum), nil
}

func day10IndicatorsMatch(result, indicatorLights []bool) bool {
	if len(result) != len(indicatorLights) {
		return false
	}
	for i, b := range result {
		if indicatorLights[i] != b {
			return false
		}
	}
	return true
}

func day10ReadMachines(rdr io.Reader) ([]day10Machine, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	machines := []day10Machine{}
	for line := range strings.SplitSeq(strings.TrimSpace(string(in)), "\n") {
		machine := day10Machine{}
		part := strings.Split(line, " ")
		if len(part) < 3 {
			return nil, fmt.Errorf("not enough parts on line: %s", line)
		}
		// parse the indicator lights in square brackets
		if part[0][0] != '[' || part[0][len(part[0])-1] != ']' {
			return nil, fmt.Errorf("indicator light missing square brackets: %s", line)
		}
		for _, c := range part[0][1 : len(part[0])-1] {
			switch c {
			case '.':
				machine.indicatorLights = append(machine.indicatorLights, false)
			case '#':
				machine.indicatorLights = append(machine.indicatorLights, true)
			default:
				return nil, fmt.Errorf("unknown indicator value %c in line: %s", c, line)
			}
		}
		// parse each button schematic in parentheses
		for _, buttonStr := range part[1 : len(part)-1] {
			buttons := []int{}
			if buttonStr[0] != '(' || buttonStr[len(buttonStr)-1] != ')' {
				return nil, fmt.Errorf("button missing parentheses on %s in line: %s", buttonStr, line)
			}
			for _, numStr := range strings.Split(buttonStr[1:len(buttonStr)-1], ",") {
				num, err := strconv.Atoi(numStr)
				if err != nil {
					return nil, fmt.Errorf("failed to parse number %s in button %s on line: %s", numStr, buttonStr, line)
				}
				buttons = append(buttons, num)
			}
			machine.buttonSchematics = append(machine.buttonSchematics, buttons)
		}
		// parse the joltage requirements in curly braces
		jPart := len(part) - 1
		if part[jPart][0] != '{' || part[jPart][len(part[jPart])-1] != '}' {
			return nil, fmt.Errorf("joltage missing curly brackets in line: %s", line)
		}
		for _, numStr := range strings.Split(part[jPart][1:len(part[jPart])-1], ",") {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse number %s in joltage %s on line: %s", numStr, part[jPart], line)
			}
			machine.joltageReq = append(machine.joltageReq, num)
		}
		machines = append(machines, machine)
	}
	return machines, nil
}

// <https://en.wikipedia.org/wiki/Gaussian_elimination>
func day10GaussianElimination(matrix [][]int, rowStart, colStart int) {
	// swap rows to get those with a non-zero col entry up top
	if rowStart >= len(matrix) || colStart >= len(matrix[rowStart]) {
		return
	}
	rowStop := rowStart - 1
	for row1 := rowStart; row1 < len(matrix); row1++ {
		if matrix[row1][colStart] != 0 {
			rowStop = row1
			continue
		}
		row2 := len(matrix) - 1
		for row2 > row1 && matrix[row2][colStart] == 0 {
			row2--
		}
		if row2 <= row1 {
			break
		}
		day10Swap(matrix, row1, row2)
		rowStop = row1
	}
	// normalize rows to all start with a positive, preferably a 1
	rowOne := -1
	for row := rowStart; row <= rowStop; row++ {
		if matrix[row][colStart] < 0 {
			day10Mult(matrix, row, -1)
		}
		// reduce to 1 if possible
		if matrix[row][colStart] > 1 {
			divAllowed := true
			for col := colStart + 1; divAllowed && col < len(matrix[row]); col++ {
				if matrix[row][col]%matrix[row][colStart] != 0 {
					divAllowed = false
				}
			}
			if divAllowed {
				div := matrix[row][colStart]
				for col := colStart; col < len(matrix[row]); col++ {
					matrix[row][col] /= div
				}
			}
		}
		// track first row starting with 1 and swap to rowStart if found
		if matrix[row][colStart] == 1 && rowOne < 0 {
			rowOne = row
		}
	}
	if rowOne >= 0 && rowOne != rowStart {
		day10Swap(matrix, rowStart, rowOne)
	}
	// reduce following rows with multiply and add
	if rowStart < rowStop && matrix[rowStart][colStart] >= 1 {
		for row2 := rowStart + 1; row2 <= rowStop; row2++ {
			multi1, multi2 := matrix[rowStart][colStart], matrix[row2][colStart]*-1
			day10Mult(matrix, row2, multi1)
			day10Add(matrix, rowStart, row2, multi2)
		}
	}
	// recurse on remaining rows/columns
	day10GaussianElimination(matrix, min(rowStart, rowStop)+1, colStart+1)
	if rowStart > rowStop || matrix[rowStart][colStart] == 0 {
		return
	}
	// attempt to reduce with division if row does not start with a 1, but do not go into fractions
	if matrix[rowStart][colStart] > 1 {
		divAllowed := true
		for col := colStart + 1; divAllowed && col < len(matrix[rowStart]); col++ {
			if matrix[rowStart][col]%matrix[rowStart][colStart] != 0 {
				divAllowed = false
			}
		}
		if divAllowed {
			div := matrix[rowStart][colStart]
			for col := colStart; col < len(matrix[rowStart]); col++ {
				matrix[rowStart][col] /= div
			}
		}
	}
	// use back-substitution to remove this field from above rows
	if matrix[rowStart][colStart] >= 1 {
		for row2 := rowStart - 1; row2 >= 0; row2-- {
			if matrix[row2][colStart] != 0 {
				multi1, multi2 := matrix[rowStart][colStart], matrix[row2][colStart]*-1
				day10Mult(matrix, row2, multi1)
				day10Add(matrix, rowStart, row2, multi2)
			}
		}
	}
}

func day10Swap(matrix [][]int, a, b int) {
	for col := range len(matrix[a]) {
		matrix[a][col], matrix[b][col] = matrix[b][col], matrix[a][col]
	}
}

func day10Mult(matrix [][]int, row, val int) {
	if val == 1 {
		return
	}
	for col := range len(matrix[row]) {
		matrix[row][col] *= val
	}
}

func day10Add(matrix [][]int, src, tgt, mult int) {
	for col := range len(matrix[src]) {
		matrix[tgt][col] += (matrix[src][col] * mult)
	}
}
