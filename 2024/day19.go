package main

import (
	"fmt"
	"io"
	"strings"
)

func init() {
	registerDay("19a", day19a)
	registerDay("19b", day19b)
}

func day19a(args []string, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	inSplit := strings.Split(string(in), "\n")
	if len(inSplit) < 3 || strings.TrimSpace(inSplit[1]) != "" {
		return "", fmt.Errorf("invalid input")
	}
	patterns := strings.Split(inSplit[0], ", ")

	sum := 0
	for num, goal := range inSplit[2:] {
		if strings.TrimSpace(goal) == "" {
			continue
		}
		fmt.Printf("searching for %d - %s using\n", num, goal)
		list := []int{0}
		for len(list) > 0 {
			// increment over list options until a match is found or options are exhausted
			result := ""
			for _, i := range list {
				result += patterns[i]
			}
			if result == goal {
				sum++
				break
			}
			if strings.HasPrefix(goal, result) {
				// increment by adding another pattern to the list
				list = append(list, 0)
			} else {
				// increment last entry in the list
				for len(list) > 0 {
					list[len(list)-1]++
					if list[len(list)-1] < len(patterns) {
						// stop on next valid entry
						break
					}
					list = list[:len(list)-1]
				}
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day19b(args []string, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	inSplit := strings.Split(string(in), "\n")
	if len(inSplit) < 3 || strings.TrimSpace(inSplit[1]) != "" {
		return "", fmt.Errorf("invalid input")
	}
	patterns := strings.Split(inSplit[0], ", ")
	tailCount := map[string]int{}

	sum := 0
	for num, goal := range inSplit[2:] {
		if strings.TrimSpace(goal) == "" {
			continue
		}
		fmt.Printf("searching for %d - %s\n", num, goal)
		for chars := range goal {
			tail := goal[len(goal)-chars-1:]
			if _, ok := tailCount[tail]; ok {
				continue
			}
			// for a given tail, count the number of combinations that can generate the tail
			tc := 0
			for _, p := range patterns {
				if !strings.HasPrefix(tail, p) {
					continue
				}
				if tail == p {
					tc++
				} else {
					tc += tailCount[tail[len(p):]]
				}
			}
			tailCount[tail] = tc
		}
		sum += tailCount[goal]
	}
	return fmt.Sprintf("%d", sum), nil
}
