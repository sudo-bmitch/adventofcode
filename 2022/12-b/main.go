package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	heights := []string{}

	in := bufio.NewScanner(os.Stdin)
	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		if len(heights) > 0 && len(line) != len(heights[0]) {
			fmt.Fprintf(os.Stderr, "uneven line lengths on %s\n", line)
			return
		}
		// TODO: consider validating line contents (a-z, S and E)
		heights = append(heights, line)
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}
	if len(heights) == 0 {
		fmt.Fprintf(os.Stderr, "missing input")
		return
	}

	// track how many steps to each position
	rows := len(heights)
	cols := len(heights[0])
	paths := [][]int{}
	paths = make([][]int, rows)
	for i := range paths {
		paths[i] = make([]int, cols)
	}

	// build up step map
	foundStep := true
	endSteps := 0
	var step int
step:
	for step = 1; foundStep; step++ {
		foundStep = false
		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				if (step == 1 && heights[row][col] == 'a') || (step > 1 && paths[row][col] == step-1) {
					foundStep = tryStep(step, heights[row][col], row-1, col, &heights, &paths) || foundStep
					foundStep = tryStep(step, heights[row][col], row+1, col, &heights, &paths) || foundStep
					foundStep = tryStep(step, heights[row][col], row, col-1, &heights, &paths) || foundStep
					foundStep = tryStep(step, heights[row][col], row, col+1, &heights, &paths) || foundStep
				}
				if heights[row][col] == 'E' && paths[row][col] != 0 {
					endSteps = paths[row][col]
					break step
				}
			}
		}
	}
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			fmt.Printf("%03d ", paths[row][col])
		}
		fmt.Printf("\n")
	}

	fmt.Printf("Result: %d\n", endSteps)
}

func tryStep(step int, from byte, row int, col int, heights *[]string, paths *[][]int) bool {
	// validate attempt: not from E, out of range, or to an already stepped on path
	if from == 'E' || row < 0 || col < 0 || row >= len(*heights) || col >= len((*heights)[0]) || (*paths)[row][col] > 0 {
		return false
	}
	// don't return to the start
	if rune((*heights)[row][col]) == 'S' {
		return false
	}
	max := rune('b')
	if step > 1 {
		max = rune(from) + 1
	}
	// if step is too high and not the end from y or z
	if ((*heights)[row][col] != 'E' && rune((*heights)[row][col]) > max) || ((*heights)[row][col] == 'E' && from < 'y') {
		return false
	}
	// mark step in paths
	(*paths)[row][col] = step
	return true
}
