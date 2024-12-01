package main

import (
	"fmt"
	"io"
	"sort"

	"github.com/sudo-bmitch/adventofcode/pkg/parse"
)

func init() {
	registerDay("01a", day01a)
	registerDay("01b", day01b)
}

func day01a(args []string, rdr io.Reader) (string, error) {
	// parse the two lists
	listA := []int{}
	listB := []int{}
	for split := range parse.MustNumSlice(rdr) {
		if len(split) != 2 {
			return "", fmt.Errorf("line does not have 2 entries: %v", split)
		}
		listA = append(listA, split[0])
		listB = append(listB, split[1])
	}
	// sort list
	sort.Ints(listA)
	sort.Ints(listB)
	// sum the differences
	sum := 0
	for i := range listA {
		a := listA[i]
		b := listB[i]
		if a > b {
			sum += a - b
		} else {
			sum += b - a
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day01b(args []string, rdr io.Reader) (string, error) {
	// parse the two lists
	listA := []int{}
	countB := map[int]int{}
	for split := range parse.MustNumSlice(rdr) {
		if len(split) != 2 {
			return "", fmt.Errorf("line does not have 2 entries: %v", split)
		}
		listA = append(listA, split[0])
		countB[split[1]]++
	}
	// sum the similarities
	sum := 0
	for _, a := range listA {
		sum += a * countB[a]
	}
	return fmt.Sprintf("%d", sum), nil
}
