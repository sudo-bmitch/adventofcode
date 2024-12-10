package main

import (
	"fmt"
	"io"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("10a", day10a)
	registerDay("10b", day10b)
}

func day10a(args []string, rdr io.Reader) (string, error) {
	return day10Run(args, rdr, true)
}

func day10b(args []string, rdr io.Reader) (string, error) {
	return day10Run(args, rdr, false)
}

type day10Path struct {
	p grid.Pos
	d grid.Dir
}

func day10Run(_ []string, rdr io.Reader, trackVisited bool) (string, error) {
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}
	sum := 0
	for p, r := range g.Walk() {
		// find each trail head
		if r != '0' {
			continue
		}
		// for each trail head, attempt to walk each path up to find peak, only visiting each position once
		visited, err := grid.New(g.W, g.H)
		if err != nil {
			return "", err
		}
		visited.G[p.X][p.Y] = '#'
		path := []day10Path{
			{p: p, d: grid.North},
		}
		for len(path) > 0 {
			// try moving
			head := path[len(path)-1]
			try := head.p.MoveD(head.d)
			if g.ValidPos(try) && g.G[head.p.X][head.p.Y]+1 == g.G[try.X][try.Y] && (!trackVisited || visited.G[try.X][try.Y] != '#') {
				// if successful, add to visited and path, inc sum if peak
				path = append(path, day10Path{p: try, d: grid.North})
				visited.G[try.X][try.Y] = '#'
				if g.G[try.X][try.Y] == '9' {
					sum++
				}
			} else {
				// else turn direction or pop from path if last direction
				for len(path) > 0 && path[len(path)-1].d+1 >= grid.DirLen {
					// remove path entries that are at the last search direction
					path = path[:len(path)-1]
				}
				// if entry remains, turn
				if len(path) > 0 {
					path[len(path)-1].d++
				}
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}
