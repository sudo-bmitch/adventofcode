package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("07a", day07a)
	registerDay("07b", day07b)
}

func day07a(args []string, rdr io.Reader) (string, error) {
	debug := false
	sum := 0
	// read a grid of runes
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}
	if debug {
		fmt.Fprint(os.Stderr, g.String())
	}
	queue := map[grid.Pos]bool{}
	// add the start
	for p, r := range g.Walk() {
		if r == 'S' {
			if debug {
				fmt.Fprintf(os.Stderr, "starting point: %s\n", p.String())
			}
			queue[p] = true
			break
		}
	}
	for len(queue) > 0 {
		// shift each queue entry down one, splitting the beam where required
		nextQueue := map[grid.Pos]bool{}
		for p := range queue {
			g.G[p.Row][p.Col] = '|' // track where we've been for debugging
			np := p.MoveD(grid.South)
			if !g.ValidPos(np) {
				continue
			}
			if g.G[np.Row][np.Col] == '^' {
				// split the beam
				sum++
				if npW := np.MoveD(grid.West); g.ValidPos(npW) {
					nextQueue[npW] = true
				}
				if npE := np.MoveD(grid.East); g.ValidPos(npE) {
					nextQueue[npE] = true
				}
			} else {
				nextQueue[np] = true
			}
		}
		queue = nextQueue
	}
	if debug {
		fmt.Fprint(os.Stderr, g.String())
	}
	return fmt.Sprintf("%d", sum), nil
}

func day07b(args []string, rdr io.Reader) (string, error) {
	debug := false
	sum := 0
	// read a grid of runes
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}
	if debug {
		fmt.Fprint(os.Stderr, g.String())
	}
	// track the number of beams at a location
	queue := map[grid.Pos]int{}
	// add the start
	for p, r := range g.Walk() {
		if r == 'S' {
			if debug {
				fmt.Fprintf(os.Stderr, "starting point: %s\n", p.String())
			}
			queue[p] = 1
			break
		}
	}
	for len(queue) > 0 {
		// shift each queue entry down one, splitting the beam where required
		nextQueue := map[grid.Pos]int{}
		for p, v := range queue {
			g.G[p.Row][p.Col] = '|' // track where we've been for debugging
			np := p.MoveD(grid.South)
			if !g.ValidPos(np) {
				sum += v // count the number of beams that made it to the bottom of the grid
				continue
			}
			if g.G[np.Row][np.Col] == '^' {
				// split the beam
				if npW := np.MoveD(grid.West); g.ValidPos(npW) {
					nextQueue[npW] += v
				}
				if npE := np.MoveD(grid.East); g.ValidPos(npE) {
					nextQueue[npE] += v
				}
			} else {
				nextQueue[np] += v
			}
		}
		queue = nextQueue
	}
	if debug {
		fmt.Fprint(os.Stderr, g.String())
	}
	return fmt.Sprintf("%d", sum), nil
}
