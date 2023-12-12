package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type day12Record struct {
	springs     []rune
	brokenSpans []int
}

func day12a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	lineNum := 0
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		r, err := day12Parse(line)
		if err != nil {
			return "", err
		}
		c := r.combinationsFast()
		lineNum++
		fmt.Fprintf(os.Stderr, "line %d: found %d combinations for %s\n", lineNum, c, line)
		sum += c
	}
	return fmt.Sprintf("%d", sum), nil
}

// note, this is a fail
func day12bParallel(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	lineNum := 0
	results := make(chan int)
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		r, err := day12Parse(line)
		if err != nil {
			return "", err
		}
		// unfold map, repeat 5 times
		nSprings := len(r.springs)
		nBroken := len(r.brokenSpans)
		for i := 1; i < 5; i++ {
			r.springs = append(r.springs, '?')
			r.springs = append(r.springs, r.springs[:nSprings]...)
			r.brokenSpans = append(r.brokenSpans, r.brokenSpans[:nBroken]...)
		}
		go func(r day12Record, lineNum int) {
			c := r.combinationsSlow()
			results <- c
		}(r, lineNum)
		lineNum++
	}
	for n := 0; n < lineNum; n++ {
		sum += <-results
		fmt.Fprintf(os.Stderr, "sum %d, remain %d\n", sum, lineNum-n-1)
	}
	return fmt.Sprintf("%d", sum), nil
}

func day12b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	lineNum := 0
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		r, err := day12Parse(line)
		if err != nil {
			return "", err
		}
		// unfold map, repeat 5 times
		nSprings := len(r.springs)
		nBroken := len(r.brokenSpans)
		for i := 1; i < 5; i++ {
			r.springs = append(r.springs, '?')
			r.springs = append(r.springs, r.springs[:nSprings]...)
			r.brokenSpans = append(r.brokenSpans, r.brokenSpans[:nBroken]...)
		}
		c := r.combinationsFast()
		lineNum++
		fmt.Fprintf(os.Stderr, "line %d: found %d combinations for %s\n", lineNum, c, line)
		sum += c
	}
	return fmt.Sprintf("%d", sum), nil
}

func day12Parse(line string) (day12Record, error) {
	r := day12Record{
		brokenSpans: []int{},
	}
	split := strings.Split(line, " ")
	if len(split) != 2 {
		return r, fmt.Errorf("missing space on line %s", line)
	}
	r.springs = []rune(split[0])
	for _, s := range r.springs {
		switch s {
		case '.', '#', '?':
			// good
		default:
			return r, fmt.Errorf("unknown spring value %c on line %s", s, line)
		}
	}
	intS := strings.Split(split[1], ",")
	for _, is := range intS {
		i, err := strconv.Atoi(is)
		if err != nil {
			return r, fmt.Errorf("failed to convert to int %s on line %s: %w", is, line, err)
		}
		r.brokenSpans = append(r.brokenSpans, i)
	}
	return r, nil
}

// this was the initial attempt, and I killed it after running for over 8 hours on the second part
func (r day12Record) combinationsSlow() int {
	// try each permutation, counting every possible permutation
	// permutations can be tracked by counting working springs (dots)
	// leading/trailing dot may be 0 length, others must be at least 1
	// sum of all working springs = count of springs - sum of broken springs
	// TODO: solutions for [0, n+1, x, x] = some combination of previous solutions? [0, n, x>1, x]
	// track all the right shifted solutions
	working := make([]int, len(r.brokenSpans))
	for i := 1; i < len(working); i++ {
		working[i] = 1
	}
	broken := 0
	for _, b := range r.brokenSpans {
		broken += b
	}
	total := len(r.springs)
	remain := total - broken - (len(r.brokenSpans) - 1)
	end := remain
	sum := 0
	// for each combination
	for {
		// compare spans to list of springs
		match := true
		pos := 0
		badSpan := 0
		for i := range r.brokenSpans {
			for n := 0; match && n < working[i]; n++ {
				if r.springs[pos] == '#' {
					match = false
				}
				pos++
			}
			for n := 0; match && n < r.brokenSpans[i]; n++ {
				if r.springs[pos] == '.' {
					match = false
				}
				pos++
			}
			if !match {
				badSpan = i
				break
			}
		}
		for n := pos; match && n < len(r.springs); n++ {
			if r.springs[n] == '#' {
				match = false
				badSpan = len(working) - 1
			}
		}
		// fmt.Fprintf(os.Stderr, "try working %v, remain %d, match %t, badSpan %d\n", working, remain, match, badSpan)
		// count if matching combination
		if match {
			sum++
			// fmt.Fprintf(os.Stderr, "match %d: %v\n", sum, working)
		}
		// increment or end
		if working[0] >= end {
			break
		}
		if !match {
			// reset spans to the right to 1, and increment bad span
			for i := badSpan + 1; i < len(working); i++ {
				remain += working[i] - 1
				working[i] = 1
			}
			if remain > 0 {
				working[badSpan]++
				remain--
			} else {
				for i := badSpan; i > 0; i-- {
					if working[i] > 1 {
						remain = working[i] - 2 // move all but 2 into remain
						working[i-1]++          // move one to next(previous) entry
						working[i] = 1          // and leave one for current entry
						break                   // found the match
					}
				}
			}
		} else if remain > 0 {
			working[len(working)-1]++
			remain--
		} else {
			for i := len(working) - 1; i > 0; i-- {
				if working[i] > 1 {
					remain = working[i] - 2 // move all but 2 into remain
					working[i-1]++          // move one to next(previous) entry
					working[i] = 1          // and leave one for current entry
					break                   // found the match
				}
			}
		}
	}
	return sum
}

func (r day12Record) combinationsFast() int {
	// build a possibilities table, range is the possible offsets from the left most position for each span
	// first entry counts 1 for each location it can be in, 0 if there's a working machine to the left or inside the span
	// next entry is sum of possible previous positions (where span is not on a working entry and gap between spans has no broken entries)
	// repeat to last span
	// count entries from last span where no broken entries found to the right

	broken := 0
	for _, b := range r.brokenSpans {
		broken += b
	}
	wiggle := len(r.springs) - broken - (len(r.brokenSpans) - 1) // how much space the spans can move
	possible := make([][]int, len(r.brokenSpans))
	for i := range possible {
		possible[i] = make([]int, wiggle+1)
	}

	for s := range r.brokenSpans {
		offset := s // offset is count of spaces between spans + broken units in each span
		for p := 0; p < s; p++ {
			offset += r.brokenSpans[p]
		}
		for w := 0; w <= wiggle; w++ {
			ok, min := r.possible(s, offset, w)
			if !ok {
				continue
			}
			if s == 0 {
				possible[s][w] = 1
				continue
			}
			for i := min; i <= w; i++ {
				possible[s][w] += possible[s-1][i]
			}
		}
	}

	sum := 0
	for end := 0; end <= wiggle; end++ {
		if end > 0 && r.springs[len(r.springs)-end] == '#' {
			break
		}
		sum += possible[len(r.brokenSpans)-1][wiggle-end]
	}

	// for s := range r.brokenSpans {
	// 	fmt.Fprintf(os.Stderr, "span %d: ", s)
	// 	for w := 0; w <= wiggle; w++ {
	// 		fmt.Fprintf(os.Stderr, "%d ", possible[s][w])
	// 	}
	// 	fmt.Fprintf(os.Stderr, "\n")
	// }

	return sum
}

func (r day12Record) possible(s, offset, w int) (bool, int) {
	// check the current broken span for working entries
	for i := 0; i < r.brokenSpans[s]; i++ {
		if r.springs[offset+w+i] == '.' {
			return false, 0
		}
	}
	// check first span for any broken entries to the left
	if s == 0 {
		for i := 0; i < w; i++ {
			if r.springs[i] == '#' {
				return false, 0
			}
		}
		return true, w
	}
	// else check other spans for gap, then smallest previous wiggle possible
	for min := w; min >= 0; min-- {
		if r.springs[offset+min-1] == '#' {
			if min < w {
				return true, min + 1
			} else {
				// broken entry in the gap
				return false, 0
			}
		}
	}
	// no min found
	return true, 0
}
