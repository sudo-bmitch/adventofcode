package main

import (
	"bufio"
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
	in := bufio.NewScanner(rdr)
	rules, err := day05ReadRules(in)
	if err != nil {
		return "", err
	}
	// parse each update
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pagesSplit := strings.Split(line, ",")
		if len(pagesSplit)%2 == 0 {
			return "", fmt.Errorf("even number of pages on line %s", line)
		}
		center := 0
		pageMap := map[int]int{}
		for i, s := range pagesSplit {
			page, err := strconv.Atoi(s)
			if err != nil {
				return "", fmt.Errorf("unable to parse %s[%d] from %s: %w", s, i, line, err)
			}
			if i == (len(pagesSplit)-1)/2 {
				center = page
			}
			pageMap[page] = i
		}
		// check against each rule
		correct := true
		for _, rule := range rules {
			p1, ok1 := pageMap[rule.page]
			p2, ok2 := pageMap[rule.later]
			if ok1 && ok2 && p1 > p2 {
				correct = false
				break
			}
		}
		if correct {
			sum += center
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day05b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	rules, err := day05ReadRules(in)
	if err != nil {
		return "", err
	}
	// parse each update
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pagesSplit := strings.Split(line, ",")
		if len(pagesSplit)%2 == 0 {
			return "", fmt.Errorf("even number of pages on line %s", line)
		}
		pages := make([]int, len(pagesSplit))
		pageMap := map[int]int{}
		for i, s := range pagesSplit {
			page, err := strconv.Atoi(s)
			if err != nil {
				return "", fmt.Errorf("unable to parse %s[%d] from %s: %w", s, i, line, err)
			}
			pages[i] = page
			pageMap[page] = i
		}
		// find all matching rules
		correct := true
		matchedRules := []day05Rule{}
		for _, rule := range rules {
			p1, ok1 := pageMap[rule.page]
			p2, ok2 := pageMap[rule.later]
			if ok1 && ok2 && p1 > p2 {
				correct = false
			}
			if ok1 && ok2 {
				matchedRules = append(matchedRules, rule)
			}
		}
		if correct {
			continue // skip the correct entries
		}
		fixedPages := make([]int, 0, len(pagesSplit))
		remainingPages := make([]int, len(pages))
		copy(remainingPages, pages)
		for len(remainingPages) > 0 {
			lenPages := len(remainingPages)
			for i := len(remainingPages) - 1; i >= 0; i-- {
				// go through remaining pages to find one that is not listed as "later" in the matchedRules
				page := remainingPages[i]
				later := false
				for _, rule := range matchedRules {
					if rule.later == page {
						later = true
						break
					}
				}
				if later {
					continue
				}
				// for the added page, remove all matchedRules for that page
				for r := len(matchedRules) - 1; r >= 0; r-- {
					if matchedRules[r].page == page {
						if r == len(matchedRules)-1 {
							matchedRules = matchedRules[:r]
						} else {
							matchedRules = append(matchedRules[:r], matchedRules[r+1:]...)
						}
					}
				}
				fixedPages = append(fixedPages, page)
				if i == len(remainingPages)-1 {
					remainingPages = remainingPages[:i]
				} else {
					remainingPages = append(remainingPages[:i], remainingPages[i+1:]...)
				}
			}
			if lenPages == len(remainingPages) {
				return "", fmt.Errorf("unable to find a page to add from %v out of %v with rules %v", remainingPages, pages, matchedRules)
			}
		}
		// Potential efficiency tip: I only need to get the center page, there's no need to track the fixedPages list.
		// Just drop entries from remainingPages until reaching the center and add that to the sum
		sum += fixedPages[(len(pages)-1)/2]
	}
	return fmt.Sprintf("%d", sum), nil
}

func day05ReadRules(in *bufio.Scanner) ([]day05Rule, error) {
	rules := []day05Rule{}
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			break // finished reading rules
		}
		ruleSplit := strings.SplitN(line, "|", 2)
		if len(ruleSplit) != 2 {
			return rules, fmt.Errorf("misformatted rule: %s", line)
		}
		page, err := strconv.Atoi(ruleSplit[0])
		if err != nil {
			return rules, fmt.Errorf("failed to parse page: %w", err)
		}
		later, err := strconv.Atoi(ruleSplit[1])
		if err != nil {
			return rules, fmt.Errorf("failed to parse later: %w", err)
		}
		rules = append(rules, day05Rule{page: page, later: later})
	}
	return rules, nil
}

type day05Rule struct {
	page, later int
}
