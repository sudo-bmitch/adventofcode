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
	const crtLen = 240
	const crtLine = 40

	valCur := 1
	valNext := 0
	cycleCur := 1
	cycleNext := 0
	crt := [crtLen]rune{}
	crtLast := 0

	in := bufio.NewScanner(os.Stdin)
	// read moves
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// split command
		cmd := strings.Split(line, " ")
		switch cmd[0] {
		case "noop":
			cycleNext = cycleCur + 1
			valNext = valCur
		case "addx":
			if len(cmd) != 2 {
				fmt.Fprintf(os.Stderr, "addx needs one arg: %s\n", line)
				return
			}
			i, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "addx input must be an int: %s: %v\n", line, err)
				return
			}
			cycleNext = cycleCur + 2
			valNext = valCur + i
		default:
			fmt.Fprintf(os.Stderr, "unknown command: %s\n", line)
			return
		}
		// handle crt
		for crtLast < cycleNext-1 && crtLast < crtLen {
			start := valCur - 1
			end := valCur + 1
			linePos := crtLast % crtLine
			if linePos >= start && linePos <= end {
				crt[crtLast] = '#'
			} else {
				crt[crtLast] = '.'
			}
			crtLast++
		}
		// update to next
		valCur = valNext
		cycleCur = cycleNext
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}
	for crtLast < cycleNext-1 && crtLast < crtLen {
		start := valCur - 1
		end := valCur + 1
		linePos := crtLast % crtLine
		if linePos >= start && linePos <= end {
			crt[crtLast] = '#'
		} else {
			crt[crtLast] = '.'
		}
		crtLast++
	}

	// show result
	for crtCur := 0; crtCur < crtLen; crtCur++ {
		linePos := crtCur % crtLine
		fmt.Printf("%c", crt[crtCur])
		if linePos == crtLine-1 {
			fmt.Printf("\n")
		}
	}

}
