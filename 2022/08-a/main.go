package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	heights := [][]int{}
	visible := [][]bool{}
	in := bufio.NewScanner(os.Stdin)
	// read in height grid
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		row, err := sToISlice(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed converting line to ints: %v\n", err)
			return
		}
		if len(heights) > 0 && len(heights[0]) != len(row) {
			fmt.Fprintf(os.Stderr, "uneven row lengths: %d to %d on %s\n", len(heights[0]), len(row), line)
			return
		}
		heights = append(heights, row)
		visible = append(visible, make([]bool, len(row)))
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}
	// compute visibilities from 4 directions
	// left to right
	for i := 0; i < len(heights); i++ {
		min := -1
		for j := 0; j < len(heights[i]); j++ {
			if heights[i][j] > min {
				min = heights[i][j]
				visible[i][j] = true
			}
		}
	}
	// top to bottom
	for j := 0; j < len(heights[0]); j++ {
		min := -1
		for i := 0; i < len(heights); i++ {
			if heights[i][j] > min {
				min = heights[i][j]
				visible[i][j] = true
			}
		}
	}
	// right to left
	for i := 0; i < len(heights); i++ {
		min := -1
		for j := len(heights[i]) - 1; j >= 0; j-- {
			if heights[i][j] > min {
				min = heights[i][j]
				visible[i][j] = true
			}
		}
	}
	// bottom to top
	for j := 0; j < len(heights[0]); j++ {
		min := -1
		for i := len(heights) - 1; i >= 0; i-- {
			if heights[i][j] > min {
				min = heights[i][j]
				visible[i][j] = true
			}
		}
	}
	// count the visible
	count := 0
	for i := 0; i < len(visible); i++ {
		for j := 0; j < len(visible[i]); j++ {
			if visible[i][j] {
				fmt.Printf("T")
				count++
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("Result: %d\n", count)
}

func sToISlice(str string) ([]int, error) {
	result := make([]int, len(str))
	var err error
	for i, c := range str {
		result[i], err = strconv.Atoi(string(c))
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
