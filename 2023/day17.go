package main

import (
	"fmt"
	"io"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

type day17Best [][][grid.DirLen]int

var day17Turns = [grid.DirLen][]grid.Dir{
	grid.North: {grid.East, grid.West},
	grid.East:  {grid.North, grid.South},
	grid.South: {grid.East, grid.West},
	grid.West:  {grid.North, grid.South},
}

type day17PathEnd struct {
	pos  grid.Pos
	dir  grid.Dir
	heat int
}

func day17a(args []string, rdr io.Reader) (string, error) {
	g, err := day17Parse(rdr)
	if err != nil {
		return "", err
	}

	best := make(day17Best, g.H)
	for x := range best {
		best[x] = make([][grid.DirLen]int, g.W)
	}

	ends := []day17PathEnd{
		{pos: grid.Pos{X: 0, Y: 0}, dir: grid.South},
		{pos: grid.Pos{X: 0, Y: 0}, dir: grid.East},
	}
	curBest := 0
	endBest := -1
	for {
		if len(ends) == 0 {
			break
		}
		e := 0
		nextBest := -1
		for e < len(ends) {
			if endBest >= 0 && ends[e].heat > endBest {
				// better ending already found, delete
				ends[e] = ends[len(ends)-1]
				ends = ends[:len(ends)-1]
				continue
			}
			if ends[e].heat > curBest {
				if nextBest < 0 || ends[e].heat < nextBest {
					nextBest = ends[e].heat
				}
				// walk a diff path first
				e++
				continue
			}
			// get the 2 turns and make up to 6 new ends for the 3 steps each
			for _, t := range day17Turns[ends[e].dir] {
				curHeat := ends[e].heat
				for steps := 1; steps <= 3; steps++ {
					dPos := grid.DirPos[t]
					dPos.X *= steps
					dPos.Y *= steps
					tryPos := ends[e].pos.MoveP(dPos)
					if !g.ValidPos(tryPos) {
						break // off the map
					}
					curHeat += int(g.G[tryPos.X][tryPos.Y] - '0')
					if endBest > 0 && curHeat > endBest {
						break // path gets too hot
					}
					if best[tryPos.X][tryPos.Y][t] > 0 && best[tryPos.X][tryPos.Y][t] <= curHeat {
						continue // another path already reached here going this way
					}
					// found a new best, add the end
					best[tryPos.X][tryPos.Y][t] = curHeat
					ends = append(ends, day17PathEnd{pos: tryPos, dir: t, heat: curHeat})
				}
			}
			// delete the cur end after it has been processed
			ends[e] = ends[len(ends)-1]
			ends = ends[:len(ends)-1]
		}
		curBest = nextBest
		// check for endBest
		for d := grid.Dir(0); d < grid.DirLen; d++ {
			if best[g.H-1][g.W-1][d] > 0 && (endBest < 0 || best[g.H-1][g.W-1][d] < endBest) {
				endBest = best[g.H-1][g.W-1][d]
			}
		}
		// fmt.Fprintf(os.Stderr, "Status: curBest=%d, ends=%v\n", curBest, ends)
	}

	return fmt.Sprintf("%d", endBest), nil
}

func day17b(args []string, rdr io.Reader) (string, error) {
	g, err := day17Parse(rdr)
	if err != nil {
		return "", err
	}

	best := make(day17Best, g.H)
	for x := range best {
		best[x] = make([][grid.DirLen]int, g.W)
	}

	ends := []day17PathEnd{
		{pos: grid.Pos{X: 0, Y: 0}, dir: grid.South},
		{pos: grid.Pos{X: 0, Y: 0}, dir: grid.East},
	}
	curBest := 0
	endBest := -1
	for {
		if len(ends) == 0 {
			break
		}
		e := 0
		nextBest := -1
		for e < len(ends) {
			if endBest >= 0 && ends[e].heat > endBest {
				// better ending already found, delete
				ends[e] = ends[len(ends)-1]
				ends = ends[:len(ends)-1]
				continue
			}
			if ends[e].heat > curBest {
				if nextBest < 0 || ends[e].heat < nextBest {
					nextBest = ends[e].heat
				}
				// walk a diff path first
				e++
				continue
			}
			// get the 2 turns and add the possible results from each
			for _, t := range day17Turns[ends[e].dir] {
				curHeat := ends[e].heat
				// take first three steps, add heat from each
				for steps := 1; steps <= 3; steps++ {
					dPos := grid.DirPos[t]
					dPos.X *= steps
					dPos.Y *= steps
					tryPos := ends[e].pos.MoveP(dPos)
					if !g.ValidPos(tryPos) {
						break // off the map
					}
					curHeat += int(g.G[tryPos.X][tryPos.Y] - '0')
				}
				for steps := 4; steps <= 10; steps++ {
					dPos := grid.DirPos[t]
					dPos.X *= steps
					dPos.Y *= steps
					tryPos := ends[e].pos.MoveP(dPos)
					if !g.ValidPos(tryPos) {
						break // off the map
					}
					curHeat += int(g.G[tryPos.X][tryPos.Y] - '0')
					if endBest > 0 && curHeat > endBest {
						break // path gets too hot
					}
					if best[tryPos.X][tryPos.Y][t] > 0 && best[tryPos.X][tryPos.Y][t] <= curHeat {
						continue // another path already reached here going this way
					}
					// found a new best, add the end
					best[tryPos.X][tryPos.Y][t] = curHeat
					ends = append(ends, day17PathEnd{pos: tryPos, dir: t, heat: curHeat})
				}
			}
			// delete the cur end after it has been processed
			ends[e] = ends[len(ends)-1]
			ends = ends[:len(ends)-1]
		}
		curBest = nextBest
		// check for endBest
		for d := grid.Dir(0); d < grid.DirLen; d++ {
			if best[g.H-1][g.W-1][d] > 0 && (endBest < 0 || best[g.H-1][g.W-1][d] < endBest) {
				endBest = best[g.H-1][g.W-1][d]
			}
		}
		// fmt.Fprintf(os.Stderr, "Status: curBest=%d, ends=%v\n", curBest, ends)
	}

	return fmt.Sprintf("%d", endBest), nil
}

func day17Parse(rdr io.Reader) (grid.Grid, error) {
	g, err := grid.FromReader(rdr)
	if err != nil {
		return g, err
	}
	for _, row := range g.G {
		for _, c := range row {
			if c < '0' || c > '9' {
				return g, fmt.Errorf("invalid grid contents on row %s, character \"%c\"", string(row), c)
			}
		}
	}
	return g, nil
}
