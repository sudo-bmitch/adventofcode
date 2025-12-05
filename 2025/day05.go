package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func init() {
	registerDay("05a", day05a)
	registerDay("05b", day05b)
}

func day05a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	inSplit := strings.SplitN(string(in), "\n\n", 2)
	if len(inSplit) != 2 {
		return "", fmt.Errorf("did not find 2 sections of input")
	}
	ranges := [][2]int{}
	for line := range strings.SplitSeq(inSplit[0], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		split := strings.SplitN(line, "-", 2)
		if len(split) != 2 {
			return "", fmt.Errorf("line does not contain a range: %s", line)
		}
		first, err := strconv.Atoi(split[0])
		if err != nil {
			return "", fmt.Errorf("failed to parse first part of range (%s) from line %s: %w", split[0], line, err)
		}
		second, err := strconv.Atoi(split[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse second part of range (%s) from line %s: %w", split[1], line, err)
		}
		ranges = append(ranges, [2]int{first, second})
	}
	for line := range strings.SplitSeq(inSplit[1], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		val, err := strconv.Atoi(line)
		if err != nil {
			return "", fmt.Errorf("failed to parse value in line %s: %w", line, err)
		}
		for _, r := range ranges {
			if r[0] <= val && val <= r[1] {
				sum++
				break
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day05b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	inSplit := strings.SplitN(string(in), "\n\n", 2)
	if len(inSplit) != 2 {
		return "", fmt.Errorf("did not find 2 sections of input")
	}
	ranges := [][2]int{}
	low, high := -1, -1
	for line := range strings.SplitSeq(inSplit[0], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		split := strings.SplitN(line, "-", 2)
		if len(split) != 2 {
			return "", fmt.Errorf("line does not contain a range: %s", line)
		}
		first, err := strconv.Atoi(split[0])
		if err != nil {
			return "", fmt.Errorf("failed to parse first part of range (%s) from line %s: %w", split[0], line, err)
		}
		second, err := strconv.Atoi(split[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse second part of range (%s) from line %s: %w", split[1], line, err)
		}
		ranges = append(ranges, [2]int{first, second})
		if low == -1 {
			low, high = first, second
		}
		low = min(low, first)
		high = max(high, second)
	}
	// no need to iterate over every number in the range, just count the number of remaining items in the current range
	// and then jump the search to the next best range or to where the current pointer made it, whichever is greater.
	cur := low
	for {
		last := cur
		jump := high + 1
		for _, r := range ranges {
			if r[0] <= cur && cur <= r[1] {
				sum += r[1] - cur + 1
				cur = r[1] + 1
			}
			if last < r[0] && r[0] < jump {
				jump = r[0]
			}
		}
		cur = max(cur, jump)
		if cur > high {
			break
		}
	}
	return fmt.Sprintf("%d", sum), nil
}
