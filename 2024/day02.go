package main

import (
	"fmt"
	"io"

	"github.com/sudo-bmitch/adventofcode/pkg/parse"
)

func init() {
	registerDay("02a", day02a)
	registerDay("02b", day02b)
}

func day02a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	for split := range parse.MustNumSlice(rdr) {
		if day02Check(split) {
			sum++
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day02b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	for split := range parse.MustNumSlice(rdr) {
		if day02Check(split) {
			sum++
			continue
		}
		// retry with first entry removed
		dropOne := make([]int, len(split)-1)
		copy(dropOne, split[1:])
		if day02Check(dropOne) {
			sum++
			continue
		}
		// iterate with remaining entries removed
		for i := 0; i < len(split)-1; i++ {
			dropOne[i] = split[i]
			if day02Check(dropOne) {
				sum++
				break
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day02Check(list []int) bool {
	if len(list) < 2 {
		return false
	}
	maxChange := 3
	isInc := list[0] < list[1]
	for i := 1; i < len(list); i++ {
		if isInc && (list[i-1] >= list[i] || list[i-1]+maxChange < list[i]) {
			return false
		}
		if !isInc && (list[i-1] <= list[i] || list[i-1]-maxChange > list[i]) {
			return false
		}
	}
	return true
}
