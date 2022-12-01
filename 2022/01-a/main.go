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
	in := bufio.NewScanner(os.Stdin)
	max := 0
	sum := 0
	for in.Scan() {
		line := in.Text()

		line = strings.TrimSpace(line)
		if line == "" {
			if sum > max {
				max = sum
			}
			sum = 0
			continue
		}

		cur, err := strconv.Atoi(string(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse %s: %v", line, err)
			return
		}

		sum += cur
	}
	if sum > max {
		max = sum
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v", err)
	}
	fmt.Printf("max is %d\n", max)

}
