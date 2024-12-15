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
			g.G[p.X][p.Y] = '.'
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
		switch g.G[next.X][next.Y] {
		default:
			return "", fmt.Errorf("unknown map contents %c at position %s", g.G[next.X][next.Y], next.String())
		case '#':
			// block moving into a wall
		case '.':
			// allow moving to empty space
			guard = next
		case 'O':
			// try pushing boxes
			for i := 1; ; i++ {
				check := guard.MoveP(grid.Pos{X: move.X * i, Y: move.Y * i})
				if !g.ValidPos(check) || g.G[check.X][check.Y] == '#' {
					// hit the wall, no where to push boxes
					break
				}
				if g.G[check.X][check.Y] == '.' {
					// found a spot, shift boxes, move guard
					g.G[check.X][check.Y] = 'O'
					g.G[next.X][next.Y] = '.'
					guard = next
					break
				}
			}
		}
	}

	sum := 0
	for p, r := range g.Walk() {
		if r == 'O' {
			sum += (p.X * 100) + p.Y
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
	for p, r := range gOrig.Walk() {
		px := p.X
		py0 := p.Y * 2
		py1 := py0 + 1
		switch r {
		case '#':
			g.G[px][py0], g.G[px][py1] = '#', '#'
		case 'O':
			g.G[px][py0], g.G[px][py1] = '[', ']'
		case '.':
			g.G[px][py0], g.G[px][py1] = '.', '.'
		case '@':
			guard.X, guard.Y = px, py0
			g.G[px][py0], g.G[px][py1] = '.', '.'
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
		switch g.G[next.X][next.Y] {
		default:
			return "", fmt.Errorf("unknown map contents %c at position %s", g.G[next.X][next.Y], next.String())
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
			moveG.G[next.X][next.Y] = g.G[next.X][next.Y]
			if move.X != 0 {
				if g.G[next.X][next.Y] == '[' {
					moveG.G[next.X][next.Y+1] = ']'
				}
				if g.G[next.X][next.Y] == ']' {
					moveG.G[next.X][next.Y-1] = '['
				}
			}
			// fail if object moves to end of map or wall, succeed if nothing was moving in last scanned row
			wall := false
			done := false
			moveInner := grid.Pos{X: move.Y * move.Y, Y: move.X * move.X}
			scanOuter := grid.Pos{X: next.X * move.X * move.X, Y: next.Y * move.Y * move.Y}
			for g.ValidPos(scanOuter) && !wall && !done {
				scanInner := scanOuter
				done = true
				for g.ValidPos(scanInner) && !wall {
					if moveG.G[scanInner.X][scanInner.Y] == '[' || moveG.G[scanInner.X][scanInner.Y] == ']' {
						scanNext := scanInner.MoveP(move)
						if !g.ValidPos(scanNext) || g.G[scanNext.X][scanNext.Y] == '#' {
							wall = true
							break
						}
						if g.G[scanNext.X][scanNext.Y] != '.' {
							// mark the object as moving
							moveG.G[scanNext.X][scanNext.Y] = g.G[scanNext.X][scanNext.Y]
							// for vertical movements, also include other half of the box
							if move.X != 0 {
								if g.G[scanNext.X][scanNext.Y] == '[' {
									moveG.G[scanNext.X][scanNext.Y+1] = ']'
								}
								if g.G[scanNext.X][scanNext.Y] == ']' {
									moveG.G[scanNext.X][scanNext.Y-1] = '['
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
			moveBack := grid.Pos{X: move.X * -1, Y: move.Y * -1}
			done = false
			scanOuter = scanOuter.MoveP(moveBack)
			for g.ValidPos(scanOuter) && !done {
				scanInner := scanOuter
				done = true
				for g.ValidPos(scanInner) {
					if moveG.G[scanInner.X][scanInner.Y] != 0 {
						moveNext := scanInner.MoveP(move)
						g.G[moveNext.X][moveNext.Y] = moveG.G[scanInner.X][scanInner.Y]
						g.G[scanInner.X][scanInner.Y] = '.'
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
			sum += (p.X * 100) + p.Y
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
