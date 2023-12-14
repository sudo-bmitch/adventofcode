package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

type day14Grid struct {
	g    [][]rune
	w, h int
}

func day14a(args []string, rdr io.Reader) (string, error) {
	g, err := day14Parse(rdr)
	if err != nil {
		return "", err
	}
	g.shiftNorth()
	sum := 0
	for x := 0; x < g.h; x++ {
		for y := 0; y < g.w; y++ {
			if g.g[x][y] == 'O' {
				sum += (g.h - x)
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day14b(args []string, rdr io.Reader) (string, error) {
	g, err := day14Parse(rdr)
	if err != nil {
		return "", err
	}
	cycle := 0
	records := map[string]int{}
	goal := 1000000000
	for cycle < goal {
		g.spin()
		cycle++
		h := g.hash()
		if prev, ok := records[h]; ok {
			diff := cycle - prev
			toGoal := goal - cycle
			if diff < toGoal {
				add := diff * int(toGoal/diff)
				fmt.Fprintf(os.Stderr, "found loop from %d to %d, adding %d to get to cycle %d\n", prev, cycle, add, cycle+add)
				cycle += add
			}
		} else {
			records[h] = cycle
		}
	}
	sum := 0
	for x := 0; x < g.h; x++ {
		for y := 0; y < g.w; y++ {
			if g.g[x][y] == 'O' {
				sum += (g.h - x)
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day14Parse(rdr io.Reader) (day14Grid, error) {
	g := day14Grid{}
	lines, err := io.ReadAll(rdr)
	if err != nil {
		return g, err
	}
	rows := strings.Split(string(lines), "\n")
	for _, r := range rows {
		if r == "" {
			continue
		}
		if g.w == 0 {
			g.w = len(r)
		}
		if len(r) != g.w {
			return g, fmt.Errorf("mixed width in grid not supported: %s", r)
		}
		g.g = append(g.g, []rune(r))
	}
	g.h = len(g.g)
	return g, nil
}

func (g day14Grid) print(w io.Writer) {
	for _, r := range g.g {
		fmt.Fprintf(w, "%s\n", string(r))
	}
}

func (g day14Grid) hash() string {
	h := sha256.New()
	g.print(h)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (g *day14Grid) spin() {
	g.shiftNorth()
	g.shiftWest()
	g.shiftSouth()
	g.shiftEast()
}

func (g *day14Grid) shiftNorth() {
	for x := 1; x < g.h; x++ {
		for y := 0; y < g.w; y++ {
			if g.g[x][y] == 'O' {
				tgtX := x
				for tgtX > 0 && g.g[tgtX-1][y] == '.' {
					tgtX--
				}
				g.g[x][y] = '.'
				g.g[tgtX][y] = 'O'
			}
		}
	}
}

func (g *day14Grid) shiftWest() {
	for x := 0; x < g.h; x++ {
		for y := 1; y < g.w; y++ {
			if g.g[x][y] == 'O' {
				tgtY := y
				for tgtY > 0 && g.g[x][tgtY-1] == '.' {
					tgtY--
				}
				g.g[x][y] = '.'
				g.g[x][tgtY] = 'O'
			}
		}
	}
}

func (g *day14Grid) shiftSouth() {
	for x := g.h - 2; x >= 0; x-- {
		for y := 0; y < g.w; y++ {
			if g.g[x][y] == 'O' {
				tgtX := x
				for tgtX < g.h-1 && g.g[tgtX+1][y] == '.' {
					tgtX++
				}
				g.g[x][y] = '.'
				g.g[tgtX][y] = 'O'
			}
		}
	}
}

func (g *day14Grid) shiftEast() {
	for x := 0; x < g.h; x++ {
		for y := g.w - 2; y >= 0; y-- {
			if g.g[x][y] == 'O' {
				tgtY := y
				for tgtY < g.w-1 && g.g[x][tgtY+1] == '.' {
					tgtY++
				}
				g.g[x][y] = '.'
				g.g[x][tgtY] = 'O'
			}
		}
	}
}
