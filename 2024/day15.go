package main

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("15a", day15a)
	registerDay("15b", day15b)
}

func day15a(args []string, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	inSplit := strings.SplitN(string(in), "\n\n", 2)
	if len(inSplit) != 2 {
		return "", fmt.Errorf("input is missing instructions")
	}
	g, err := grid.FromString(inSplit[0])
	if err != nil {
		return "", err
	}
	// find the guard's location
	guard := grid.Pos{}
	for p, r := range g.Walk() {
		if r == '@' {
			guard = p
			g.G[p.Row][p.Col] = '.'
			break
		}
	}
	// apply move instructions
	for _, m := range inSplit[1] {
		if unicode.IsSpace(m) {
			continue
		}
		moveD, ok := day15MapMoves[m]
		if !ok {
			return "", fmt.Errorf("unknown move: %c", m)
		}
		move := grid.DirPos[moveD]
		next := guard.MoveP(move)
		// block guard from moving off the map
		if !g.ValidPos(next) {
			continue
		}
		// switch based on next position
		switch g.G[next.Row][next.Col] {
		default:
			return "", fmt.Errorf("unknown map contents %c at position %s", g.G[next.Row][next.Col], next.String())
		case '#':
			// block moving into a wall
		case '.':
			// allow moving to empty space
			guard = next
		case 'O':
			// try pushing boxes
			for i := 1; ; i++ {
				check := guard.MoveP(grid.Pos{Row: move.Row * i, Col: move.Col * i})
				if !g.ValidPos(check) || g.G[check.Row][check.Col] == '#' {
					// hit the wall, no where to push boxes
					break
				}
				if g.G[check.Row][check.Col] == '.' {
					// found a spot, shift boxes, move guard
					g.G[check.Row][check.Col] = 'O'
					g.G[next.Row][next.Col] = '.'
					guard = next
					break
				}
			}
		}
	}

	sum := 0
	for p, v := range g.Walk() {
		if v == 'O' {
			sum += (p.Row * 100) + p.Col
		}
	}

	return fmt.Sprintf("%d", sum), nil
}

func day15b(args []string, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	inSplit := strings.SplitN(string(in), "\n\n", 2)
	if len(inSplit) != 2 {
		return "", fmt.Errorf("input is missing instructions")
	}
	gOrig, err := grid.FromString(inSplit[0])
	if err != nil {
		return "", err
	}
	// double the width and track guard's initial location
	g, err := grid.New(gOrig.W*2, gOrig.H)
	if err != nil {
		return "", err
	}
	guard := grid.Pos{}
	for p, v := range gOrig.Walk() {
		pr := p.Row
		pc0 := p.Col * 2
		pc1 := pc0 + 1
		switch v {
		case '#':
			g.G[pr][pc0], g.G[pr][pc1] = '#', '#'
		case 'O':
			g.G[pr][pc0], g.G[pr][pc1] = '[', ']'
		case '.':
			g.G[pr][pc0], g.G[pr][pc1] = '.', '.'
		case '@':
			guard.Row, guard.Col = pr, pc0
			g.G[pr][pc0], g.G[pr][pc1] = '.', '.'
		}
	}
	for _, m := range inSplit[1] {
		if unicode.IsSpace(m) {
			continue
		}
		moveD, ok := day15MapMoves[m]
		if !ok {
			return "", fmt.Errorf("unknown move: %c", m)
		}
		move := grid.DirPos[moveD]
		next := guard.MoveP(move)
		// block guard from moving off the map
		if !g.ValidPos(next) {
			continue
		}
		// switch based on next position
		switch g.G[next.Row][next.Col] {
		default:
			return "", fmt.Errorf("unknown map contents %c at position %s", g.G[next.Row][next.Col], next.String())
		case '#':
			// block moving into a wall
		case '.':
			// allow moving to empty space
			guard = next
		case '[', ']':
			// We do things this screwed up way not because they are easy, but because we thought they would be easy.
			// try pushing boxes
			// track position of every box that is moving, use a grid to dedup
			moveG, err := grid.New(g.W, g.H)
			if err != nil {
				return "", err
			}
			moveG.G[next.Row][next.Col] = g.G[next.Row][next.Col]
			if move.Row != 0 {
				if g.G[next.Row][next.Col] == '[' {
					moveG.G[next.Row][next.Col+1] = ']'
				}
				if g.G[next.Row][next.Col] == ']' {
					moveG.G[next.Row][next.Col-1] = '['
				}
			}
			// fail if object moves to end of map or wall, succeed if nothing was moving in last scanned row
			wall := false
			done := false
			moveInner := grid.Pos{Row: move.Col * move.Col, Col: move.Row * move.Row}
			scanOuter := grid.Pos{Row: next.Row * move.Row * move.Row, Col: next.Col * move.Col * move.Col}
			for g.ValidPos(scanOuter) && !wall && !done {
				scanInner := scanOuter
				done = true
				for g.ValidPos(scanInner) && !wall {
					if moveG.G[scanInner.Row][scanInner.Col] == '[' || moveG.G[scanInner.Row][scanInner.Col] == ']' {
						scanNext := scanInner.MoveP(move)
						if !g.ValidPos(scanNext) || g.G[scanNext.Row][scanNext.Col] == '#' {
							wall = true
							break
						}
						if g.G[scanNext.Row][scanNext.Col] != '.' {
							// mark the object as moving
							moveG.G[scanNext.Row][scanNext.Col] = g.G[scanNext.Row][scanNext.Col]
							// for vertical movements, also include other half of the box
							if move.Row != 0 {
								if g.G[scanNext.Row][scanNext.Col] == '[' {
									moveG.G[scanNext.Row][scanNext.Col+1] = ']'
								}
								if g.G[scanNext.Row][scanNext.Col] == ']' {
									moveG.G[scanNext.Row][scanNext.Col-1] = '['
								}
							}
							done = false
						}
					}
					scanInner = scanInner.MoveP(moveInner)
				}
				scanOuter = scanOuter.MoveP(move)
			}
			// if wall reached, guard doesn't move, stop
			if wall {
				continue
			}
			guard = next
			// start moving blocks, work from far edge back, placing spaces where object moves from
			moveBack := grid.Pos{Row: move.Row * -1, Col: move.Col * -1}
			done = false
			scanOuter = scanOuter.MoveP(moveBack)
			for g.ValidPos(scanOuter) && !done {
				scanInner := scanOuter
				done = true
				for g.ValidPos(scanInner) {
					if moveG.G[scanInner.Row][scanInner.Col] != 0 {
						moveNext := scanInner.MoveP(move)
						g.G[moveNext.Row][moveNext.Col] = moveG.G[scanInner.Row][scanInner.Col]
						g.G[scanInner.Row][scanInner.Col] = '.'
						done = false
					}
					scanInner = scanInner.MoveP(moveInner)
				}
				scanOuter = scanOuter.MoveP(moveBack)
			}
		}
	}

	sum := 0
	for p, r := range g.Walk() {
		if r == '[' {
			sum += (p.Row * 100) + p.Col
		}
	}

	return fmt.Sprintf("%d", sum), nil
}

var day15MapMoves = map[rune]grid.Dir{
	'^': grid.North,
	'>': grid.East,
	'v': grid.South,
	'<': grid.West,
}
