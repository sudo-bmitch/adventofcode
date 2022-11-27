package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	first := true
	last := 0
	count := 0
	for in.Scan() {
		line := in.Text()

		cur, err := strconv.Atoi(string(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse %s: %v", line, err)
			return
		}

		if !first && cur > last {
			count++
		}
		last = cur
		first = false
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v", err)
	}
	fmt.Printf("found %d increasing lines", count)
}
