package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("14a", day14a)
	registerDay("14b", day14b)
}

func day14a(args []string, rdr io.Reader) (string, error) {
	quadrants := [4]int{}
	seconds := 100
	// maxR, maxC := 11, 7 // test values
	maxR, maxC := 101, 103 // run values
	midR, midC := (maxR-1)/2, (maxC-1)/2
	lineRe := regexp.MustCompile(`^\s*p=(-?\d+),(-?\d+)\s+v=(-?\d+),(-?\d+)`)
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		lineSplit := lineRe.FindStringSubmatch(line)
		if len(lineSplit) != 5 {
			return "", fmt.Errorf("failed to parse line: %s", line)
		}
		pr, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse px from line: %s, %w", line, err)
		}
		pc, err := strconv.Atoi(lineSplit[2])
		if err != nil {
			return "", fmt.Errorf("failed to parse py from line: %s, %w", line, err)
		}
		vr, err := strconv.Atoi(lineSplit[3])
		if err != nil {
			return "", fmt.Errorf("failed to parse vx from line: %s, %w", line, err)
		}
		vc, err := strconv.Atoi(lineSplit[4])
		if err != nil {
			return "", fmt.Errorf("failed to parse vy from line: %s, %w", line, err)
		}
		pr = (pr + vr*seconds) % maxR
		pc = (pc + vc*seconds) % maxC
		if pr < 0 {
			pr += maxR
		}
		if pc < 0 {
			pc += maxC
		}
		fmt.Printf("result of %s: %d,%d\n", line, pr, pc)
		if pr < midR {
			if pc < midC {
				quadrants[0]++
			}
			if pc > midC {
				quadrants[1]++
			}
		}
		if pr > midR {
			if pc < midC {
				quadrants[2]++
			}
			if pc > midC {
				quadrants[3]++
			}
		}
	}

	sum := quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
	return fmt.Sprintf("%d", sum), nil
}

func day14b(args []string, rdr io.Reader) (string, error) {
	robots := []day14Robot{}
	maxR, maxC := 101, 103
	lineRe := regexp.MustCompile(`^\s*p=(-?\d+),(-?\d+)\s+v=(-?\d+),(-?\d+)`)
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		lineSplit := lineRe.FindStringSubmatch(line)
		if len(lineSplit) != 5 {
			return "", fmt.Errorf("failed to parse line: %s", line)
		}
		pr, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse px from line: %s, %w", line, err)
		}
		pc, err := strconv.Atoi(lineSplit[2])
		if err != nil {
			return "", fmt.Errorf("failed to parse py from line: %s, %w", line, err)
		}
		vr, err := strconv.Atoi(lineSplit[3])
		if err != nil {
			return "", fmt.Errorf("failed to parse vx from line: %s, %w", line, err)
		}
		vc, err := strconv.Atoi(lineSplit[4])
		if err != nil {
			return "", fmt.Errorf("failed to parse vy from line: %s, %w", line, err)
		}
		robots = append(robots, day14Robot{
			p: grid.Pos{Row: pr, Col: pc},
			v: grid.Pos{Row: vr, Col: vc},
		})
	}

	search := strings.Repeat("#", 10)
	found := false
	sum := 0
	for i := 0; !found; i++ {
		g, err := grid.New(maxR, maxC)
		if err != nil {
			return "", err
		}
		for p := range g.Walk() {
			g.G[p.Row][p.Col] = ' '
		}
		for _, r := range robots {
			pr := (r.p.Row + r.v.Row*i) % maxR
			pc := (r.p.Col + r.v.Col*i) % maxC
			if pr < 0 {
				pr += maxR
			}
			if pc < 0 {
				pc += maxC
			}
			g.G[pc][pr] = '#'
		}
		for _, row := range g.G {
			if strings.Contains(string(row), search) {
				found = true
			}
		}
		if found {
			fmt.Printf("Graph at second %d:\n%s\n", i, g.String())
			sum = i
		}
	}

	return fmt.Sprintf("%d", sum), nil
}

type day14Robot struct {
	p, v grid.Pos
}
