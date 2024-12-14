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
	// maxX, maxY := 11, 7 // test values
	maxX, maxY := 101, 103 // run values
	midX, midY := (maxX-1)/2, (maxY-1)/2
	lineRe := regexp.MustCompile(`^\s*p=(-?\d+),(-?\d+)\s+v=(-?\d+),(-?\d+)`)
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		lineSplit := lineRe.FindStringSubmatch(line)
		if len(lineSplit) != 5 {
			return "", fmt.Errorf("failed to parse line: %s", line)
		}
		px, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse px from line: %s, %w", line, err)
		}
		py, err := strconv.Atoi(lineSplit[2])
		if err != nil {
			return "", fmt.Errorf("failed to parse py from line: %s, %w", line, err)
		}
		vx, err := strconv.Atoi(lineSplit[3])
		if err != nil {
			return "", fmt.Errorf("failed to parse vx from line: %s, %w", line, err)
		}
		vy, err := strconv.Atoi(lineSplit[4])
		if err != nil {
			return "", fmt.Errorf("failed to parse vy from line: %s, %w", line, err)
		}
		px = (px + vx*seconds) % maxX
		py = (py + vy*seconds) % maxY
		if px < 0 {
			px += maxX
		}
		if py < 0 {
			py += maxY
		}
		fmt.Printf("result of %s: %d,%d\n", line, px, py)
		if px < midX {
			if py < midY {
				quadrants[0]++
			}
			if py > midY {
				quadrants[1]++
			}
		}
		if px > midX {
			if py < midY {
				quadrants[2]++
			}
			if py > midY {
				quadrants[3]++
			}
		}
	}

	sum := quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
	return fmt.Sprintf("%d", sum), nil
}

func day14b(args []string, rdr io.Reader) (string, error) {
	robots := []day14Robot{}
	maxX, maxY := 101, 103
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
		px, err := strconv.Atoi(lineSplit[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse px from line: %s, %w", line, err)
		}
		py, err := strconv.Atoi(lineSplit[2])
		if err != nil {
			return "", fmt.Errorf("failed to parse py from line: %s, %w", line, err)
		}
		vx, err := strconv.Atoi(lineSplit[3])
		if err != nil {
			return "", fmt.Errorf("failed to parse vx from line: %s, %w", line, err)
		}
		vy, err := strconv.Atoi(lineSplit[4])
		if err != nil {
			return "", fmt.Errorf("failed to parse vy from line: %s, %w", line, err)
		}
		robots = append(robots, day14Robot{
			p: grid.Pos{X: px, Y: py},
			v: grid.Pos{X: vx, Y: vy},
		})
	}

	search := strings.Repeat("#", 10)
	found := false
	sum := 0
	for i := 0; !found; i++ {
		g, err := grid.New(maxX, maxY)
		if err != nil {
			return "", err
		}
		for p := range g.Walk() {
			g.G[p.X][p.Y] = ' '
		}
		for _, r := range robots {
			px := (r.p.X + r.v.X*i) % maxX
			py := (r.p.Y + r.v.Y*i) % maxY
			if px < 0 {
				px += maxX
			}
			if py < 0 {
				py += maxY
			}
			g.G[py][px] = '#'
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
