package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func init() {
	registerDay("03a", day03a)
	registerDay("03b", day03b)
}

func day03a(args []string, rdr io.Reader) (string, error) {
	debug := false
	sum := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		val, pos := day03LargestRemaining(line, 0, 1)
		best := val
		val, _ = day03LargestRemaining(line, pos+1, 0)
		best = best*10 + val
		if debug {
			fmt.Fprintf(os.Stderr, "from line %s, adding %d\n", line, best)
		}
		sum += best
	}

	return fmt.Sprintf("%d", sum), nil
}

func day03b(args []string, rdr io.Reader) (string, error) {
	debug := false
	sum := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		best := 0
		lastPos := -1
		for remain := 12; remain > 0; remain-- {
			val, pos := day03LargestRemaining(line, lastPos+1, remain-1)
			best = best*10 + val
			lastPos = pos
		}
		if debug {
			fmt.Fprintf(os.Stderr, "from line %s, adding %d\n", line, best)
		}
		sum += best
	}

	return fmt.Sprintf("%d", sum), nil
}

func day03LargestRemaining(line string, first, remain int) (val, pos int) {
	val = -1
	pos = first
	for i, c := range line[first : len(line)-remain] {
		cur := int(c - '0')
		if cur > val {
			val = cur
			pos = first + i
		}
	}
	return val, pos
}
