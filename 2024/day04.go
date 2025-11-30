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
		{Row: 1},           // down
		{Row: -1},          // up
		{Col: 1},           // right
		{Col: -1},          // left
		{Row: 1, Col: 1},   // diag down/right
		{Row: -1, Col: 1},  // diag up/right
		{Row: 1, Col: -1},  // diag down/left
		{Row: -1, Col: -1}, // diag up/left
	}
	for r := 0; r < g.H; r++ {
		for c := 0; c < g.W; c++ {
			if g.G[r][c] != goal[0] {
				continue // skip everything that doesn't start with the first letter
			}
			for _, m := range moves {
				if !g.ValidPos(grid.Pos{Row: r + (m.Row * (len(goal) - 1)), Col: c + (m.Col * (len(goal) - 1))}) {
					continue // skip when searching a move goes outside of grid
				}
				found := true
				for i := 1; i < len(goal) && found; i++ {
					if g.G[r+(m.Row*i)][c+(m.Col*i)] != goal[i] {
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
	for r := 0; r < g.H-(goalSize-1); r++ {
		for c := 0; c < g.W-(goalSize-1); c++ {
			for _, goal := range goals {
				found := true
				for dr := 0; dr < goalSize && found; dr++ {
					for dc := 0; dc < goalSize && found; dc++ {
						if goal[dr][dc] == '.' {
							continue
						}
						if g.G[r+dr][c+dc] != goal[dr][dc] {
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
