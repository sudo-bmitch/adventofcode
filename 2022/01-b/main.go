package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	sum := 0
	sums := []int{}
	topN := 0
	count := 3
	for in.Scan() {
		line := in.Text()

		line = strings.TrimSpace(line)
		if line == "" {
			sums = append(sums, sum)
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
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v", err)
	}
	sums = append(sums, sum)
	sort.Sort(sort.Reverse(sort.IntSlice(sums)))
	for _, cur := range sums[:count] {
		fmt.Printf("adding %d\n", cur)
		topN += cur
	}
	fmt.Printf("top %d is %d\n", count, topN)

}
