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
	total := 0
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		r0Lo, r0Hi, r1Lo, r1Hi, err := parsePairs(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse ranges %s: %v\n", line, err)
		}
		if !(r0Lo > r1Hi || r1Lo > r0Hi) {
			total++
		}
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v", err)
	}
	fmt.Printf("Total: %d\n", total)
}

func parsePairs(line string) (int, int, int, int, error) {
	elfs := strings.SplitN(line, ",", 2)
	if len(elfs) < 2 {
		return 0, 0, 0, 0, fmt.Errorf("missing comma")
	}
	range0 := strings.SplitN(elfs[0], "-", 2)
	if len(range0) < 2 {
		return 0, 0, 0, 0, fmt.Errorf("missing dash in first range")
	}
	range1 := strings.SplitN(elfs[1], "-", 2)
	if len(range1) < 2 {
		return 0, 0, 0, 0, fmt.Errorf("missing dash in second range")
	}
	range0Lo, err := strconv.Atoi(range0[0])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to parse int %s: %w", range0[0], err)
	}
	range0Hi, err := strconv.Atoi(range0[1])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to parse int %s: %w", range0[1], err)
	}
	range1Lo, err := strconv.Atoi(range1[0])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to parse int %s: %w", range1[0], err)
	}
	range1Hi, err := strconv.Atoi(range1[1])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to parse int %s: %w", range1[1], err)
	}
	if range0Lo > range0Hi || range1Lo > range1Hi {
		return 0, 0, 0, 0, fmt.Errorf("invalid ranges, lo greater than hi value")
	}
	return range0Lo, range0Hi, range1Lo, range1Hi, nil
}
