package main

import (
	"fmt"
	"io"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("04a", day04a)
	registerDay("04b", day04b)
}

func day04a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	goal := []rune("XMAS")
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}
	moves := []grid.Pos{
		{X: 1},         // down
		{X: -1},        // up
		{Y: 1},         // right
		{Y: -1},        // left
		{X: 1, Y: 1},   // diag down/right
		{X: -1, Y: 1},  // diag up/right
		{X: 1, Y: -1},  // diag down/left
		{X: -1, Y: -1}, // diag up/left
	}
	for x := 0; x < g.H; x++ {
		for y := 0; y < g.W; y++ {
			if g.G[x][y] != goal[0] {
				continue // skip everything that doesn't start with the first letter
			}
			for _, m := range moves {
				if !g.ValidPos(grid.Pos{X: x + (m.X * (len(goal) - 1)), Y: y + (m.Y * (len(goal) - 1))}) {
					continue // skip when searching a move goes outside of grid
				}
				found := true
				for i := 1; i < len(goal) && found; i++ {
					if g.G[x+(m.X*i)][y+(m.Y*i)] != goal[i] {
						found = false
					}
				}
				if found {
					sum++
				}
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day04b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}
	goals := [][][]rune{
		{
			[]rune("M.S"),
			[]rune(".A."),
			[]rune("M.S"),
		},
		{
			[]rune("S.M"),
			[]rune(".A."),
			[]rune("S.M"),
		},
		{
			[]rune("M.M"),
			[]rune(".A."),
			[]rune("S.S"),
		},
		{
			[]rune("S.S"),
			[]rune(".A."),
			[]rune("M.M"),
		},
	}
	goalSize := 3
	for x := 0; x < g.H-(goalSize-1); x++ {
		for y := 0; y < g.W-(goalSize-1); y++ {
			for _, goal := range goals {
				found := true
				for dx := 0; dx < goalSize && found; dx++ {
					for dy := 0; dy < goalSize && found; dy++ {
						if goal[dx][dy] == '.' {
							continue
						}
						if g.G[x+dx][y+dy] != goal[dx][dy] {
							found = false
						}
					}
				}
				if found {
					sum++
				}
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}
