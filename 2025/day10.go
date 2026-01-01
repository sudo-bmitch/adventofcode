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
	out := make(chan int)

	for i, machine := range machines {
		// this should be solved with liner algebra, but I wasn't smart enough to implement that in Go.
		// check each machine concurrently to offset just how inefficient this is
		go func() {
			best := math.MaxInt
			// track number of button presses in an array, starting with all buttons at zero
			buttonPresses := make([]int, len(machine.buttonSchematics))
			for {
				result := make([]int, len(machine.joltageReq))
				// compute the result based on the current button presses
				for buttonCur, buttonCount := range buttonPresses {
					for _, counter := range machine.buttonSchematics[buttonCur] {
						result[counter] += buttonCount
					}
				}
				// check the result vs the desired joltage, note any low counters
				high, low := false, false
				lowCounters := make([]bool, len(machine.joltageReq))
				for ji, val := range result {
					if val > machine.joltageReq[ji] {
						high = true
					}
					if val < machine.joltageReq[ji] {
						low = true
						lowCounters[ji] = true
					}
				}
				if !low && !high {
					// match found, compute number of presses, and update best count if better
					count := 0
					for _, n := range buttonPresses {
						count += n
					}
					if count < best {
						if debug {
							fmt.Fprintf(os.Stderr, "machine %d, best count now %d, found button presses %v\n", i, count, buttonPresses)
						}
						best = count
					}
				}
				if low && !high {
					// at least one value is low and no values are high, increment by finding last button that would result in all low charges incrementing
					cur := len(buttonPresses) - 1
					for cur >= 0 {
						// unmark entries that the current button increments
						for _, counter := range machine.buttonSchematics[cur] {
							lowCounters[counter] = false
						}
						noneLow := true
						for _, b := range lowCounters {
							if b {
								noneLow = false
								break
							}
						}
						if noneLow {
							break
						}
						cur--
					}
					// increment current button and zero remaining buttons
					buttonPresses[cur]++
					for bi := cur + 1; bi < len(buttonPresses); bi++ {
						buttonPresses[bi] = 0
					}
				} else {
					// either a match or at least one high value, iterate by decreasing last non-zero button to zero and increment previous button
					// if only first button was non-zero, break
					cur := len(buttonPresses) - 1
					for cur >= 0 && buttonPresses[cur] == 0 {
						cur--
					}
					if cur == 0 {
						break
					}
					buttonPresses[cur] = 0
					buttonPresses[cur-1]++
				}
			}
			if debug {
				fmt.Fprintf(os.Stderr, "machine %d requires %d button presses\n", i, best)
			}
			out <- best
		}()
	}

	for i := range machines {
		sum += <-out
		if debug {
			fmt.Fprintf(os.Stderr, "Waiting on %d results, current total of %d machines is %d\n", len(machines)-i-1, i+1, sum)
		}
	}
	close(out)

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
