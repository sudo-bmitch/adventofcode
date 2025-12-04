package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("04a", day04a)
	registerDay("04b", day04b)
}

func day04a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	debug := false
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}
	dg, err := grid.New(g.W, g.H)
	if err != nil {
		return "", err
	}
	for p := range g.Walk() {
		if debug {
			dg.G[p.Row][p.Col] = g.G[p.Row][p.Col]
		}
		if g.G[p.Row][p.Col] != '@' {
			continue
		}
		adj := 0
		for _, diff := range day04Adjacent {
			check := p.MoveP(diff)
			if g.ValidPos(check) && g.G[check.Row][check.Col] == '@' {
				adj++
			}
		}
		if adj < 4 {
			sum++
			if debug {
				dg.G[p.Row][p.Col] = 'x'
			}
		}
	}
	if debug {
		fmt.Fprint(os.Stderr, dg.String())

	}
	return fmt.Sprintf("%d", sum), nil
}

func day04b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	debug := false
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}
	dg, err := grid.New(g.W, g.H)
	if err != nil {
		return "", err
	}
	for {
		removedOne := false
		for p := range g.Walk() {
			if debug {
				dg.G[p.Row][p.Col] = g.G[p.Row][p.Col]
			}
			if g.G[p.Row][p.Col] != '@' {
				continue
			}
			adj := 0
			for _, diff := range day04Adjacent {
				check := p.MoveP(diff)
				if g.ValidPos(check) && g.G[check.Row][check.Col] == '@' {
					adj++
				}
			}
			if adj < 4 {
				g.G[p.Row][p.Col] = 'x'
				removedOne = true
				sum++
				if debug {
					dg.G[p.Row][p.Col] = 'x'
				}
			}
		}
		if !removedOne {
			break
		}
	}
	if debug {
		fmt.Fprint(os.Stderr, dg.String())

	}
	return fmt.Sprintf("%d", sum), nil
}

var day04Adjacent = []grid.Pos{
	{Row: -1, Col: 0},
	{Row: -1, Col: 1},
	{Row: 0, Col: 1},
	{Row: 1, Col: 1},
	{Row: 1, Col: 0},
	{Row: 1, Col: -1},
	{Row: 0, Col: -1},
	{Row: -1, Col: -1},
}
