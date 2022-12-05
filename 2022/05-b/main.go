package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	stacks, err := parseStacks(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse stacks: %v", err)
		return
	}
	fmt.Printf("starting state:\n")
	printStacks(stacks)
	for in.Scan() {
		line := in.Text()
		// line = strings.TrimSpace(line) // do not trim line since leading whitespace is important
		err = parseMove(line, &stacks)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse move: %v", err)
			return
		}
		fmt.Printf("After %s:\n", line)
		printStacks(stacks)
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v", err)
		return
	}
	fmt.Printf("Result: ")
	for i, stack := range stacks {
		if len(stack) == 0 {
			fmt.Fprintf(os.Stderr, "empty stack: %d\n", i)
		}
		fmt.Printf("%c", stack[len(stack)-1])
	}
}

func parseStacks(in *bufio.Scanner) ([][]rune, error) {
	// read up to empty line
	stackLines := []string{}
	for in.Scan() {
		line := in.Text()
		// line = strings.TrimSpace(line) // do not trim space since leading whitespace is important
		if line == "" {
			break
		}
		stackLines = append(stackLines, line)
	}
	// read in count of stacks from end, error if there are things other than numbers in order
	indexLine := strings.TrimSpace(stackLines[len(stackLines)-1])
	indexStrings := strings.Split(indexLine, "   ")
	for i, s := range indexStrings {
		if s != fmt.Sprintf("%d", i+1) {
			return nil, fmt.Errorf("index %d is %s and should be %d", i, s, i+1)
		}
	}
	stackCount := len(indexStrings)
	stacks := make([][]rune, stackCount)
	for i := 0; i < stackCount; i++ {
		// start each stack as empty, and lines*2 is a rough guess of the starting capacity
		stacks[i] = make([]rune, 0, len(stackLines)*2)
	}
	// read each stack entry from bottom up, error if not a single character in brackets
	for lineNum := len(stackLines) - 2; lineNum >= 0; lineNum-- {
		line := stackLines[lineNum]
		for stackNum := 0; stackNum < stackCount; stackNum++ {
			// append onto each stack, error if stack doesn't have n-1 items already, skip empty columns
			charOffset := 1 + (4 * stackNum)
			if charOffset > len(line) {
				continue // empty end of line
			}
			c := line[charOffset]
			if c == byte(' ') {
				continue
			}
			if len(stacks[stackNum]) < len(stackLines)-2-lineNum {
				return nil, fmt.Errorf("possible hole in stack, reading line %d, stack %d, cur len %d", lineNum, stackNum, len(stacks[stackNum]))
			}
			stacks[stackNum] = append(stacks[stackNum], rune(c))
		}
	}
	return stacks, nil
}

var reMove = regexp.MustCompile(`^move (\d+) from (\d+) to (\d+)$`)

func parseMove(line string, stacks *[][]rune) error {
	// parse "move $count from $start to $end"
	matches := reMove.FindSubmatch([]byte(line))
	if len(matches) < 4 {
		return fmt.Errorf("could not parse %s: %v", line, matches)
	}
	count, err := strconv.Atoi(string(matches[1]))
	if err != nil {
		return fmt.Errorf("failed to convert count %s: %w", matches[1], err)
	}
	start, err := strconv.Atoi(string(matches[2]))
	if err != nil {
		return fmt.Errorf("failed to convert start %s: %w", matches[2], err)
	}
	end, err := strconv.Atoi(string(matches[3]))
	if err != nil {
		return fmt.Errorf("failed to convert end %s: %w", matches[3], err)
	}
	// offset from 0
	start = start - 1
	end = end - 1

	// error if $start stack doesn't have enough items
	if start < 0 || end < 0 || start > len(*stacks)-1 || end > len(*stacks)-1 || len((*stacks)[start]) < count {
		return fmt.Errorf("invalid start/end for stacks: count %d, start %d, end %d, stacks: %v", count, start, end, *stacks)
	}
	pop := (*stacks)[start][len((*stacks)[start])-count:]
	(*stacks)[end] = append((*stacks)[end], pop...)
	(*stacks)[start] = (*stacks)[start][:len((*stacks)[start])-count]
	return nil
}

func printStacks(stacks [][]rune) {
	for i, stack := range stacks {
		fmt.Printf("%d: ", i)
		for _, c := range stack {
			fmt.Printf("%c", c)
		}
		fmt.Printf("\n")
	}
}
