package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func init() {
	registerDay("01a", day01a)
	registerDay("01b", day01b)
}

func day01a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	cur := 50
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		dir := 0
		switch line[0] {
		case 'L':
			dir = -1
		case 'R':
			dir = 1
		default:
			return "", fmt.Errorf("unknown line: %s", line)
		}
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", fmt.Errorf("failed to parse number in %s: %w", line, err)
		}
		cur = ((cur + (dir * count)) + 100) % 100
		if cur == 0 {
			sum++
		}
		// fmt.Fprintf(os.Stderr, "sum = %d, cur = %d, line parsed = %s\n", sum, cur, line)
	}

	return fmt.Sprintf("%d", sum), nil
}

func day01b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	cur := 50
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		dir := 0
		switch line[0] {
		case 'L':
			dir = -1
		case 'R':
			dir = 1
		default:
			return "", fmt.Errorf("unknown line: %s", line)
		}
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", fmt.Errorf("failed to parse number in %s: %w", line, err)
		}
		if count >= 100 {
			sum += int(count / 100)
			count = count % 100
		}
		if cur == 0 || count == 0 {
			// noop
		} else if dir > 0 && cur+count >= 100 {
			sum++
		} else if dir < 0 && cur-count <= 0 {
			sum++
		}
		cur = ((cur + (dir * count)) + 100) % 100
		// fmt.Fprintf(os.Stderr, "sum = %d, cur = %d, line parsed = %s\n", sum, cur, line)
	}
	return fmt.Sprintf("%d", sum), nil
}
