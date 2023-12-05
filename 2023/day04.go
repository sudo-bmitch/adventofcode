package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type card struct {
	i    int
	win  []int
	have []int
}

func day04a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		c, err := day04Parse(line)
		if err != nil {
			return "", err
		}
		matches := day04Matches(c.win, c.have)
		if matches > 0 {
			sum += (1 << (matches - 1))
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day04b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	copies := map[int]int{} // a rolling limited length array/slice would be more efficient, but this is faster to write
	n := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		c, err := day04Parse(line)
		if err != nil {
			return "", err
		}
		matches := day04Matches(c.win, c.have)
		for i := 1; i <= matches; i++ {
			copies[n+i] += 1 + copies[n]
		}
		sum += 1 + copies[n]
		n++
	}
	return fmt.Sprintf("%d", sum), nil
}

func day04Parse(line string) (card, error) {
	c := card{}
	colon := strings.Split(line, ": ")
	if len(colon) != 2 || colon[0][:5] != "Card " {
		return c, fmt.Errorf("failed to parse card and colon from %s", line)
	}
	ci, err := strconv.Atoi(strings.Trim(colon[0][5:], " "))
	if err != nil {
		return c, fmt.Errorf("failed to parse card number from %s on line %s", colon[0][5:], line)
	}
	c.i = ci
	winHave := strings.Split(colon[1], " | ")
	if len(winHave) != 2 {
		return c, fmt.Errorf("failed to parse pipe from %s", colon[1])
	}
	winNums, err := parseNumListBySpace(winHave[0])
	if err != nil {
		return c, fmt.Errorf("failed to parse win numbers from %s: %w", winHave[0], err)
	}
	haveNums, err := parseNumListBySpace(winHave[1])
	if err != nil {
		return c, fmt.Errorf("failed to parse have numbers from %s: %w", winHave[1], err)
	}
	c.win = winNums
	c.have = haveNums
	return c, nil
}

func day04Matches(win, have []int) int {
	count := 0
	for _, w := range win {
		for _, h := range have {
			if w == h {
				count++
				break
			}
		}
	}
	return count
}
