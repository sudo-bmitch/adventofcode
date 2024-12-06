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
	for x := 0; x < g.H && !found; x++ {
		for y := 0; y < g.W && !found; y++ {
			if g.G[x][y] == '^' {
				g.G[x][y] = '.'
				curPos.X = x
				curPos.Y = y
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
		if g.G[nextPos.X][nextPos.Y] == '.' {
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
	for x := 0; x < g.H && !found; x++ {
		for y := 0; y < g.W && !found; y++ {
			if g.G[x][y] == '^' {
				g.G[x][y] = '.'
				startPos.X = x
				startPos.Y = y
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
		if g.G[nextPos.X][nextPos.Y] == '.' {
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
		x, y := testLoc.X, testLoc.Y
		// test with new obstruction in the location
		g.G[x][y] = '#'
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
			if g.G[nextPos.X][nextPos.Y] == '.' {
				curPos = nextPos
			} else {
				curDir = day06TurnRight[curDir]
			}
		}
		g.G[x][y] = '.' // reset map
	}
	return fmt.Sprintf("%d", sum), nil
}

var day06TurnRight = map[grid.Dir]grid.Dir{
	grid.North: grid.East,
	grid.East:  grid.South,
	grid.South: grid.West,
	grid.West:  grid.North,
}
