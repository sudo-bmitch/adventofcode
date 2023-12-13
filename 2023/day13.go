package main

import (
	"fmt"
	"io"
	"strings"
)

type day13Grid struct {
	g    [][]rune
	w, h int
}

func day13a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	gList, err := day13Parse(rdr)
	if err != nil {
		return "", err
	}
	for _, g := range gList {
		rc := g.reflectCols(0)
		rr := g.reflectRows(0)
		sum += rc + (rr * 100)
	}

	return fmt.Sprintf("%d", sum), nil
}

func day13b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	gList, err := day13Parse(rdr)
	if err != nil {
		return "", err
	}
	for _, g := range gList {
		rc := g.reflectCols(1)
		rr := g.reflectRows(1)
		sum += rc + (rr * 100)
	}

	return fmt.Sprintf("%d", sum), nil
}

func day13Parse(rdr io.Reader) ([]day13Grid, error) {
	lines, err := io.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	ret := []day13Grid{}
	grids := strings.Split(string(lines), "\n\n")
	for _, gridLines := range grids {
		if gridLines == "" {
			continue
		}
		g := day13Grid{}
		rows := strings.Split(gridLines, "\n")
		g.w = len(rows[0])
		for _, r := range rows {
			if r == "" {
				continue
			}
			if len(r) != g.w {
				return nil, fmt.Errorf("mixed width in grid not supported: %s", r)
			}
			g.g = append(g.g, []rune(r))
		}
		g.h = len(g.g)
		ret = append(ret, g)
	}
	return ret, nil
}

func (g day13Grid) reflectCols(smudge int) int {
	smudges := make([]int, g.w)
	for x := 0; x < g.h; x++ {
		for y := 1; y < g.w; y++ {
			// w=10, y=6, c=8, reflect point = 4 = 6 - (8-6+1)
			// w=10, y=2, c=3, 2-(3-2+1)
			for c := y; c < g.w && 2*y-c-1 >= 0; c++ {
				if g.g[x][c] != g.g[x][2*y-c-1] {
					smudges[y]++
				}
			}
		}
	}
	// return best
	for i := len(smudges) - 1; i > 0; i-- {
		if smudges[i] == smudge {
			return i
		}
	}
	return 0
}

func (g day13Grid) reflectRows(smudge int) int {
	smudges := make([]int, g.h)
	for y := 0; y < g.w; y++ {
		for x := 1; x < g.h; x++ {
			for c := x; c < g.h && 2*x-c-1 >= 0; c++ {
				if g.g[c][y] != g.g[2*x-c-1][y] {
					smudges[x]++
				}
			}
		}
	}
	// return best
	for i := len(smudges) - 1; i > 0; i-- {
		if smudges[i] == smudge {
			return i
		}
	}
	return 0
}
