package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

func init() {
	registerDay("03a", day03a)
	registerDay("03b", day03b)
}

func day03a(args []string, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	sum, err := day03Sum(in)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", sum), nil
}

func day03b(args []string, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	sum := 0
	re := regexp.MustCompile(`(do\(\)|don't\(\))`)
	matches := re.FindAllIndex(in, -1)
	// parse up to first match
	end := len(in)
	if len(matches) > 0 && len(matches[0]) > 0 {
		end = matches[0][0]
	}
	add, err := day03Sum(in[:end])
	if err != nil {
		return "", err
	}
	sum += add
	// parse from each "do" match to next match or end of string
	for i, match := range matches {
		if len(match) == 0 {
			return "", fmt.Errorf("match did not contain entries: %v", match)
		}
		// check if match is a "do" match
		pos := match[0]
		if string(in[pos:pos+4]) == "do()" {
			end := len(in)
			if i+1 < len(matches) && len(matches[i+1]) > 0 {
				end = matches[i+1][0]
			}
			add, err := day03Sum(in[pos+4 : end])
			if err != nil {
				return "", err
			}
			sum += add
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day03Sum(in []byte) (int, error) {
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	sum := 0
	matches := re.FindAllSubmatch(in, -1)
	for _, match := range matches {
		if len(match) != 3 {
			return 0, fmt.Errorf("match did not contain 3 entries: %v", match)
		}
		a, err := strconv.Atoi(string(match[1]))
		if err != nil {
			return 0, fmt.Errorf("failed to convert %s to int: %w", match[1], err)
		}
		b, err := strconv.Atoi(string(match[2]))
		if err != nil {
			return 0, fmt.Errorf("failed to convert %s to int: %w", match[2], err)
		}
		sum += (a * b)
	}
	return sum, nil
}
