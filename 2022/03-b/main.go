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
	for {
		// load 3 rucksacks (lines)
		if !in.Scan() {
			break
		}
		rucksack1 := in.Text()
		rucksack1 = strings.TrimSpace(rucksack1)
		if rucksack1 == "" {
			continue
		}
		if !in.Scan() {
			fmt.Fprintf(os.Stderr, "could not find a second rucksack after %s\n", rucksack1)
			return
		}
		rucksack2 := in.Text()
		rucksack2 = strings.TrimSpace(rucksack2)
		if rucksack2 == "" {
			fmt.Fprintf(os.Stderr, "second rucksack empty after %s\n", rucksack1)
			return
		}
		if !in.Scan() {
			fmt.Fprintf(os.Stderr, "could not find a third rucksack after %s and %s\n", rucksack1, rucksack2)
			return
		}
		rucksack3 := in.Text()
		rucksack3 = strings.TrimSpace(rucksack3)
		if rucksack3 == "" {
			fmt.Fprintf(os.Stderr, "third rucksack empty after %s and %s\n", rucksack1, rucksack2)
			return
		}
		// find the item in all three rucksacks
		match, err := findMatch(rucksack1, rucksack2, rucksack3)
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

func findMatch(a, b, c string) (rune, error) {
	for _, ch := range a {
		if strings.ContainsRune(b, ch) && strings.ContainsRune(c, ch) {
			return ch, nil
		}
	}
	return 0, fmt.Errorf("no matching items found in %s, %s, and %s", a, b, c)
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
