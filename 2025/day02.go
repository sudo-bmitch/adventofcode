package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func init() {
	registerDay("02a", day02a)
	registerDay("02b", day02b)
}

func day02a(args []string, rdr io.Reader) (string, error) {
	debug := false
	sum := 0
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	ranges := strings.Split(string(in), ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		parts := strings.SplitN(r, "-", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("failed to parse range %s", r)
		}
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", fmt.Errorf("failed to parse start of range %s: %w", r, err)
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse end of range %s: %w", r, err)
		}
		for i := start; i <= end; i++ {
			iStr := fmt.Sprintf("%d", i)
			if len(iStr)%2 != 0 {
				continue
			}
			mid := len(iStr) / 2
			first := iStr[:mid]
			second := iStr[mid:]
			if first == second {
				if debug {
					fmt.Fprintf(os.Stderr, "found %d in range %s\n", i, r)
				}
				sum += i
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day02b(args []string, rdr io.Reader) (string, error) {
	debug := false
	sum := 0
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	ranges := strings.Split(string(in), ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		parts := strings.SplitN(r, "-", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("failed to parse range %s", r)
		}
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", fmt.Errorf("failed to parse start of range %s: %w", r, err)
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse end of range %s: %w", r, err)
		}
		for i := start; i <= end; i++ {
			iStr := fmt.Sprintf("%d", i)
			for patLen := 1; patLen <= len(iStr)/2; patLen++ {
				if len(iStr)%patLen != 0 {
					continue
				}
				pat := iStr[:patLen]
				curMatch := true
				for p := patLen; p < len(iStr); p += patLen {
					if pat != iStr[p:p+patLen] {
						curMatch = false
						break
					}
				}
				if curMatch {
					if debug {
						fmt.Fprintf(os.Stderr, "found %d in range %s\n", i, r)
					}
					sum += i
					break
				}
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}
