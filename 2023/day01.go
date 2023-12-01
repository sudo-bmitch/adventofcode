package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

func day01a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		first := 0
		last := 0
		for i := 0; i < len(line); i++ {
			if line[i] >= '0' && line[i] <= '9' {
				first = int(line[i]) - int('0')
				break
			}
		}
		for i := len(line) - 1; i >= 0; i-- {
			if line[i] >= '0' && line[i] <= '9' {
				last = int(line[i]) - int('0')
				break
			}
		}
		sum += (first * 10) + last
	}
	return fmt.Sprintf("%d", sum), nil
}

func day01b(args []string, rdr io.Reader) (string, error) {
	debug := false
	if len(args) > 0 && args[0] == "debug" {
		debug = true
	}
	numbers := `[0-9]|zero|one|two|three|four|five|six|seven|eight|nine`
	reFirst := regexp.MustCompile(`^.*?(` + numbers + `).*$`)
	reLast := regexp.MustCompile(`^.*(` + numbers + `).*?$`)
	sum := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		matchFirst := reFirst.FindStringSubmatch(line)
		matchLast := reLast.FindStringSubmatch(line)
		if len(matchFirst) < 2 || len(matchLast) < 2 {
			return "", fmt.Errorf("failed to find match in %s", line)
		}
		first := toNumber(matchFirst[1])
		last := toNumber(matchLast[1])
		if debug {
			fmt.Fprintf(os.Stderr, "adding from %s: %d and %d\n", line, first, last)
		}
		sum += (first * 10) + last
	}
	return fmt.Sprintf("%d", sum), nil
}

func toNumber(s string) int {
	switch s {
	case "0", "zero":
		return 0
	case "1", "one":
		return 1
	case "2", "two":
		return 2
	case "3", "three":
		return 3
	case "4", "four":
		return 4
	case "5", "five":
		return 5
	case "6", "six":
		return 6
	case "7", "seven":
		return 7
	case "8", "eight":
		return 8
	case "9", "nine":
		return 9
	default:
		panic("unknown value: " + s)
	}
}
