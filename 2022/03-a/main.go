package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	lu := setupLU()
	total := 0
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// split one rucksack into 2 compartments
		split := len(line) / 2
		compA := line[:split]
		compB := line[split:]
		if len(compA) != len(compB) {
			fmt.Fprintf(os.Stderr, "Uneven compartments: %s and %s\n", compA, compB)
			return
		}
		// find the item in both compartments
		match, err := findMatch(compA, compB)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to find a match: %v\n", err)
			return
		}
		// use a lookup to find te value of the matching item
		if rune(match) > 255 || lu[rune(match)] == 0 {
			fmt.Fprintf(os.Stderr, "failed to lookup value for %c\n", match)
		}
		total += lu[rune(match)]
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v", err)
	}
	fmt.Printf("Total: %d\n", total)
}

func findMatch(a, b string) (rune, error) {
	for _, ch := range a {
		if strings.ContainsRune(b, ch) {
			return ch, nil
		}
	}
	return 0, fmt.Errorf("no matching items found in %s and %s", a, b)
}

func setupLU() [256]int {
	lu := [256]int{}
	for ch := rune('a'); ch <= rune('z'); ch++ {
		lu[ch] = int(rune(ch) - rune('a') + 1)
	}
	for ch := rune('A'); ch <= rune('Z'); ch++ {
		lu[ch] = int(rune(ch) - rune('A') + 27)
	}
	return lu
}
