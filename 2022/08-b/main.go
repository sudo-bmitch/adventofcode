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
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}
	// compute views
	bestView := 0
	for i := 0; i < len(heights); i++ {
		for j := 0; j < len(heights[i]); j++ {
			curHeight := heights[i][j]
			left := 0
			for look := j - 1; look >= 0; look-- {
				left++
				if heights[i][look] >= curHeight {
					break
				}
			}
			right := 0
			for look := j + 1; look < len(heights[i]); look++ {
				right++
				if heights[i][look] >= curHeight {
					break
				}
			}
			up := 0
			for look := i - 1; look >= 0; look-- {
				up++
				if heights[look][j] >= curHeight {
					break
				}
			}
			down := 0
			for look := i + 1; look < len(heights); look++ {
				down++
				if heights[look][j] >= curHeight {
					break
				}
			}
			view := left * right * up * down
			if view > bestView {
				bestView = view
			}
		}
	}
	fmt.Printf("Result: %d\n", bestView)
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
