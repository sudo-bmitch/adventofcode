package main

import (
	"fmt"
	"io"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("06a", day06a)
	registerDay("06b", day06b)
}

func day06a(args []string, rdr io.Reader) (string, error) {
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	// find starting position
	curPos := grid.Pos{}
	curDir := grid.North
	found := false
	for r := 0; r < g.H && !found; r++ {
		for c := 0; c < g.W && !found; c++ {
			if g.G[r][c] == '^' {
				g.G[r][c] = '.'
				curPos.Row = r
				curPos.Col = c
				found = true
			}
		}
	}
	// walk the map until leaving the edge
	visits := map[grid.Pos]bool{}
	for {
		visits[curPos] = true
		nextPos := curPos.MoveD(curDir)
		if !g.ValidPos(nextPos) {
			break
		}
		if g.G[nextPos.Row][nextPos.Col] == '.' {
			curPos = nextPos
		} else {
			curDir = day06TurnRight[curDir]
		}
	}
	sum := len(visits)
	return fmt.Sprintf("%d", sum), nil
}

func day06b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	// find starting position
	startPos := grid.Pos{}
	found := false
	for r := 0; r < g.H && !found; r++ {
		for c := 0; c < g.W && !found; c++ {
			if g.G[r][c] == '^' {
				g.G[r][c] = '.'
				startPos.Row = r
				startPos.Col = c
				found = true
			}
		}
	}
	// get a list of places the guard normally walks
	curPos := startPos
	curDir := grid.North
	visits := map[grid.Pos]bool{}
	for {
		visits[curPos] = true
		nextPos := curPos.MoveD(curDir)
		if !g.ValidPos(nextPos) {
			break
		}
		if g.G[nextPos.Row][nextPos.Col] == '.' {
			curPos = nextPos
		} else {
			curDir = day06TurnRight[curDir]
		}
	}
	// only test obstructions in locations the guard would normally walk
	for testLoc := range visits {
		if testLoc == startPos {
			continue // skip testing start position
		}
		r, c := testLoc.Row, testLoc.Col
		// test with new obstruction in the location
		g.G[r][c] = '#'
		curPos := startPos
		curDir := grid.North
		// track loops by observing where the same position and direction are repeated
		visitDir := map[grid.Pos]*[grid.DirLen]bool{}
		for {
			if visitDir[curPos] == nil {
				visitDir[curPos] = &[grid.DirLen]bool{}
			}
			if visitDir[curPos][curDir] {
				// found a loop
				sum++
				break
			}
			visitDir[curPos][curDir] = true
			nextPos := curPos.MoveD(curDir)
			if !g.ValidPos(nextPos) {
				// walked off the map
				break
			}
			if g.G[nextPos.Row][nextPos.Col] == '.' {
				curPos = nextPos
			} else {
				curDir = day06TurnRight[curDir]
			}
		}
		g.G[r][c] = '.' // reset map
	}
	return fmt.Sprintf("%d", sum), nil
}

var day06TurnRight = map[grid.Dir]grid.Dir{
	grid.North: grid.East,
	grid.East:  grid.South,
	grid.South: grid.West,
	grid.West:  grid.North,
}
