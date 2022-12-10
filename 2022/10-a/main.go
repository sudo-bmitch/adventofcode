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
	sum := 0
	valCur := 1
	valNext := 0
	timeCur := 1
	timeNext := 0
	sampleNext := 20
	sampleDiff := 40

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
			timeNext = timeCur + 1
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
			timeNext = timeCur + 2
			valNext = valCur + i
		default:
			fmt.Fprintf(os.Stderr, "unknown command: %s\n", line)
			return
		}
		// handle sample
		if timeNext > sampleNext {
			sum += valCur * sampleNext
			fmt.Printf("sample %d at %d, sum %d\n", valCur, sampleNext, sum)
			sampleNext += sampleDiff
		}
		// update to next
		valCur = valNext
		timeCur = timeNext
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}
	fmt.Printf("Result: %d\n", sum)
}
