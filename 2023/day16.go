package main

import (
	"fmt"
	"io"
	"strings"
)

type day16Grid struct {
	g    [][]rune
	w, h int
}

type day16Dir int

const (
	day16North day16Dir = iota
	day16East
	day16South
	day16West
)

type day16Pos struct {
	x, y int
}

func (p day16Pos) move(m day16Pos) day16Pos {
	return day16Pos{x: p.x + m.x, y: p.y + m.y}
}

var day16Move = map[day16Dir]day16Pos{
	day16North: {x: -1},
	day16East:  {y: 1},
	day16South: {x: 1},
	day16West:  {y: -1},
}

type day16Ray struct {
	p day16Pos
	d day16Dir
}

var day16Mirror = map[rune]map[day16Dir][]day16Dir{
	'.': {
		day16North: []day16Dir{day16North},
		day16East:  []day16Dir{day16East},
		day16South: []day16Dir{day16South},
		day16West:  []day16Dir{day16West},
	},
	'/': {
		day16North: []day16Dir{day16East},
		day16East:  []day16Dir{day16North},
		day16South: []day16Dir{day16West},
		day16West:  []day16Dir{day16South},
	},
	'\\': {
		day16North: []day16Dir{day16West},
		day16East:  []day16Dir{day16South},
		day16South: []day16Dir{day16East},
		day16West:  []day16Dir{day16North},
	},
	'-': {
		day16North: []day16Dir{day16East, day16West},
		day16East:  []day16Dir{day16East},
		day16South: []day16Dir{day16East, day16West},
		day16West:  []day16Dir{day16West},
	},
	'|': {
		day16North: []day16Dir{day16North},
		day16East:  []day16Dir{day16North, day16South},
		day16South: []day16Dir{day16South},
		day16West:  []day16Dir{day16North, day16South},
	},
}

func day16a(args []string, rdr io.Reader) (string, error) {
	g, err := day16Parse(rdr)
	if err != nil {
		return "", err
	}
	sum := g.compute(day16Ray{p: day16Pos{x: 0, y: 0}, d: day16East})
	return fmt.Sprintf("%d", sum), nil
}

func day16b(args []string, rdr io.Reader) (string, error) {
	g, err := day16Parse(rdr)
	if err != nil {
		return "", err
	}
	best := 0
	for x := 0; x < g.h; x++ {
		sum := g.compute(day16Ray{p: day16Pos{x: x, y: 0}, d: day16East})
		if sum > best {
			best = sum
		}
		sum = g.compute(day16Ray{p: day16Pos{x: x, y: g.w - 1}, d: day16West})
		if sum > best {
			best = sum
		}
	}
	for y := 0; y < g.w; y++ {
		sum := g.compute(day16Ray{p: day16Pos{x: 0, y: y}, d: day16South})
		if sum > best {
			best = sum
		}
		sum = g.compute(day16Ray{p: day16Pos{x: g.h - 1, y: y}, d: day16North})
		if sum > best {
			best = sum
		}
	}
	return fmt.Sprintf("%d", best), nil
}

func day16Parse(rdr io.Reader) (day16Grid, error) {
	g := day16Grid{
		g: [][]rune{},
	}
	in, err := io.ReadAll(rdr)
	if err != nil {
		return g, err
	}
	for _, line := range strings.Split(string(in), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if g.w == 0 {
			g.w = len(line)
		} else if g.w != len(line) {
			return g, fmt.Errorf("line length inconsistent on line %s", line)
		}
		g.g = append(g.g, []rune(line))
	}
	g.h = len(g.g)
	return g, nil
}

func (g day16Grid) compute(start day16Ray) int {
	rays := []day16Ray{start}
	// visited[x][y][dir] == bool
	visited := make([][][day16West + 1]bool, g.h)
	for x := range visited {
		visited[x] = make([][day16West + 1]bool, g.w)
	}
	visited[start.p.x][start.p.y][start.d] = true
	for len(rays) > 0 {
		nextRays := []day16Ray{}
		for _, r := range rays {
			dirs := day16Mirror[g.g[r.p.x][r.p.y]][r.d]
			for _, dir := range dirs {
				move := day16Move[dir]
				nextPos := r.p.move(move)
				// fmt.Fprintf(os.Stderr, "moved from [%d,%d x %d] to [%d,%d x %d] using move [%d,%d]\n", r.p.x, r.p.y, r.d, nextPos.x, nextPos.y, dir, move.x, move.y)
				if nextPos.x < 0 || nextPos.y < 0 || nextPos.x >= g.h || nextPos.y >= g.w {
					// fmt.Fprintf(os.Stderr, "dropped ray out of range: [%d,%d]\n", nextPos.x, nextPos.y)
					continue
				}
				if visited[nextPos.x][nextPos.y][dir] {
					// fmt.Fprintf(os.Stderr, "dropped ray already visited: [%d,%d]\n", nextPos.x, nextPos.y)
					continue
				}
				visited[nextPos.x][nextPos.y][dir] = true
				nextRays = append(nextRays, day16Ray{p: nextPos, d: dir})
			}
		}
		rays = nextRays
	}
	sum := 0
	for x := 0; x < g.h; x++ {
		for y := 0; y < g.w; y++ {
			if visited[x][y][day16North] || visited[x][y][day16East] || visited[x][y][day16South] || visited[x][y][day16West] {
				sum++
			}
		}
	}
	return sum
}
