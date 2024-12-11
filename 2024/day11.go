package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func init() {
	registerDay("11a", day11a)
	registerDay("11b", day11b)
}

func day11a(args []string, rdr io.Reader) (string, error) {
	return day11Run(rdr, 25)
}

func day11b(args []string, rdr io.Reader) (string, error) {
	return day11Run(rdr, 75)
}

func day11Run(rdr io.Reader, numBlinks int) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	split := strings.Fields(string(in))
	nums := make([]int, len(split))
	for i, s := range split {
		nums[i], err = strconv.Atoi(s)
		if err != nil {
			return "", err
		}
	}

	nextMap := map[int]*day11NumNext{}
	curNums := nums
	for i := range numBlinks {
		nextNums := make([]int, 0, len(curNums)*2)
		for _, n := range curNums {
			if nextMap[n] != nil {
				continue
			}
			next := day11Blink(n)
			nextMap[n] = &day11NumNext{next: next, maxLen: numBlinks - i}
			nextNums = append(nextNums, next...)
		}
		if len(nextNums) == 0 {
			break
		}
		curNums = nextNums
	}

	// populate the childLenList for each map entry until the stack height is reached for that number
	for i := range numBlinks {
		for n := range nextMap {
			// skip entries where maxLen is reached
			if nextMap[n].maxLen <= i {
				continue
			}
			// append the sum of childLenList of the previous step from the children
			if i == 0 {
				nextMap[n].childLenList = make([]int, numBlinks)
				nextMap[n].childLenList[0] = len(nextMap[n].next)
			} else {
				sum := 0
				for _, c := range nextMap[n].next {
					sum += nextMap[c].childLenList[i-1]
				}
				nextMap[n].childLenList[i] = sum
			}
		}
	}

	// now sum up each num at the blink count
	sum := 0
	for _, n := range nums {
		sum += nextMap[n].childLenList[numBlinks-1]
	}
	return fmt.Sprintf("%d", sum), nil
}

func day11Blink(n int) []int {
	out := []int{}
	if n == 0 {
		out = append(out, 1)
	} else if numLen := day11NumLen(n); numLen%2 == 0 {
		div := 10
		for i := 1; i < numLen/2; i++ {
			div *= 10
		}
		out = append(out, n/div, n%div)
	} else {
		out = append(out, n*2024)
	}
	return out
}

func day11NumLen(i int) int {
	cmp, n := 10, 1
	for i >= cmp {
		cmp *= 10
		n++
	}
	return n
}

type day11NumNext struct {
	next         []int
	maxLen       int
	childLenList []int
}
