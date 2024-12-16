package main

import (
	"fmt"
	"io"
	"slices"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("16a", day16a)
	registerDay("16b", day16b)
}

func day16a(args []string, rdr io.Reader) (string, error) {
	startDir := grid.East
	startPos := grid.Pos{}
	endPos := grid.Pos{}
	addStraight := 1
	addTurn := 1000

	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}

	// locate starting position
	for p, r := range g.Walk() {
		if r == 'S' {
			startPos = p
		}
		if r == 'E' {
			endPos = p
		}
	}
	if (startPos == grid.Pos{}) {
		return "", fmt.Errorf("start position not found")
	}
	if (endPos == grid.Pos{}) {
		return "", fmt.Errorf("end position not found")
	}

	bestScore := 0
	bestPath := day16Path{}
	bestAtPos := map[grid.Pos]int{}
	paths := []day16Path{
		{
			p:     []grid.Pos{startPos},
			d:     startDir,
			score: 0,
		},
	}
	for len(paths) > 0 {
		slices.SortFunc(paths, func(a, b day16Path) int {
			if a.score < b.score {
				return -1
			}
			if a.score > b.score {
				return 1
			}
			return 0
		})
		path := paths[0]
		curPos := path.p[len(path.p)-1]
		// check for reaching the end and if so, stop
		if curPos == endPos {
			bestPath = path
			bestScore = path.score
			break
		}
		straightPos := curPos.MoveD(path.d)
		leftPos := curPos.MoveD(day16Turns[path.d][0])
		rightPos := curPos.MoveD(day16Turns[path.d][1])
		// if turn is valid, append to list of paths
		if g.ValidPos(leftPos) && g.G[leftPos.X][leftPos.Y] != '#' && !path.loop(leftPos) && (bestAtPos[leftPos] == 0 || bestAtPos[leftPos] > path.score) {
			leftPath := path.clone()
			leftPath.p = append(leftPath.p, leftPos)
			leftPath.d = day16Turns[path.d][0]
			leftPath.score += addTurn + addStraight
			paths = append(paths, leftPath)
			if bs, ok := bestAtPos[leftPos]; !ok || bs > leftPath.score {
				bestAtPos[leftPos] = leftPath.score
			}
		}
		if g.ValidPos(rightPos) && g.G[rightPos.X][rightPos.Y] != '#' && !path.loop(rightPos) && (bestAtPos[rightPos] == 0 || bestAtPos[rightPos] > path.score) {
			rightPath := path.clone()
			rightPath.p = append(rightPath.p, rightPos)
			rightPath.d = day16Turns[path.d][1]
			rightPath.score += addTurn + addStraight
			paths = append(paths, rightPath)
			if bs, ok := bestAtPos[rightPos]; !ok || bs > rightPath.score {
				bestAtPos[rightPos] = rightPath.score
			}
		}
		if g.ValidPos(straightPos) && g.G[straightPos.X][straightPos.Y] != '#' && !path.loop(straightPos) && (bestAtPos[straightPos] == 0 || bestAtPos[straightPos]+addTurn > path.score) {
			// increment straight
			paths[0].p = append(paths[0].p, straightPos)
			paths[0].score += addStraight
			if bs, ok := bestAtPos[straightPos]; !ok || bs > paths[0].score {
				bestAtPos[straightPos] = paths[0].score
			}
		} else {
			// failed going straight, remove from paths
			if len(paths) == 1 {
				paths = paths[:0]
			} else {
				paths = paths[1:]
			}
		}
	}

	for _, p := range bestPath.p {
		g.G[p.X][p.Y] = '*'
	}
	fmt.Printf("best path found:\n%s\n", g.String())

	return fmt.Sprintf("%d", bestScore), nil
}

func day16b(args []string, rdr io.Reader) (string, error) {
	startDir := grid.East
	startPos := grid.Pos{}
	endPos := grid.Pos{}
	addStraight := 1
	addTurn := 1000

	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}

	// locate starting position
	for p, r := range g.Walk() {
		if r == 'S' {
			startPos = p
		}
		if r == 'E' {
			endPos = p
		}
	}
	if (startPos == grid.Pos{}) {
		return "", fmt.Errorf("start position not found")
	}
	if (endPos == grid.Pos{}) {
		return "", fmt.Errorf("end position not found")
	}

	bestScore := 0
	bestAtPos := map[grid.Pos]int{}
	paths := []day16Path{
		{
			p:     []grid.Pos{startPos},
			d:     startDir,
			score: 0,
		},
	}
	bestPosList := []grid.Pos{}
	for len(paths) > 0 {
		slices.SortFunc(paths, func(a, b day16Path) int {
			if a.score < b.score {
				return -1
			}
			if a.score > b.score {
				return 1
			}
			return 0
		})
		path := paths[0]
		curPos := path.p[len(path.p)-1]
		if bestScore > 0 && path.score > bestScore {
			// done searching paths, all remaining entries have a worse score
			break
		}
		// check for reaching the end and if so, track positions on the path and remove from list
		if curPos == endPos {
			bestPosList = append(bestPosList, path.p...)
			bestScore = path.score
			if len(paths) == 1 {
				paths = paths[:0]
			} else {
				paths = paths[1:]
			}
			continue
		}
		straightPos := curPos.MoveD(path.d)
		leftPos := curPos.MoveD(day16Turns[path.d][0])
		rightPos := curPos.MoveD(day16Turns[path.d][1])
		// if turn is valid, append to list of paths
		if g.ValidPos(leftPos) && g.G[leftPos.X][leftPos.Y] != '#' && !path.loop(leftPos) && (bestAtPos[leftPos] == 0 || bestAtPos[leftPos]+(addTurn*3) > path.score) {
			leftPath := path.clone()
			leftPath.p = append(leftPath.p, leftPos)
			leftPath.d = day16Turns[path.d][0]
			leftPath.score += addTurn + addStraight
			paths = append(paths, leftPath)
			if bs, ok := bestAtPos[leftPos]; !ok || bs > leftPath.score {
				bestAtPos[leftPos] = leftPath.score
			}
		}
		if g.ValidPos(rightPos) && g.G[rightPos.X][rightPos.Y] != '#' && !path.loop(rightPos) && (bestAtPos[rightPos] == 0 || bestAtPos[rightPos]+(addTurn*3) > path.score) {
			rightPath := path.clone()
			rightPath.p = append(rightPath.p, rightPos)
			rightPath.d = day16Turns[path.d][1]
			rightPath.score += addTurn + addStraight
			paths = append(paths, rightPath)
			if bs, ok := bestAtPos[rightPos]; !ok || bs > rightPath.score {
				bestAtPos[rightPos] = rightPath.score
			}
		}
		if g.ValidPos(straightPos) && g.G[straightPos.X][straightPos.Y] != '#' && !path.loop(straightPos) && (bestAtPos[straightPos] == 0 || bestAtPos[straightPos]+(addTurn*3) > path.score) {
			// increment straight
			paths[0].p = append(paths[0].p, straightPos)
			paths[0].score += addStraight
			if bs, ok := bestAtPos[straightPos]; !ok || bs > paths[0].score {
				bestAtPos[straightPos] = paths[0].score
			}
		} else {
			// failed going straight, remove from paths
			if len(paths) == 1 {
				paths = paths[:0]
			} else {
				paths = paths[1:]
			}
		}
	}

	countSeats := 0
	for _, p := range bestPosList {
		if g.G[p.X][p.Y] != 'O' {
			g.G[p.X][p.Y] = 'O'
			countSeats++
		}
	}
	fmt.Printf("seat map:\n%s\n", g.String())
	return fmt.Sprintf("%d", countSeats), nil
}

var day16Turns = map[grid.Dir][2]grid.Dir{
	grid.North: {grid.West, grid.East},
	grid.East:  {grid.North, grid.South},
	grid.South: {grid.East, grid.West},
	grid.West:  {grid.South, grid.North},
}

type day16Path struct {
	p     []grid.Pos
	d     grid.Dir
	score int
}

func (path day16Path) clone() day16Path {
	c := day16Path{
		p:     make([]grid.Pos, len(path.p)),
		d:     path.d,
		score: path.score,
	}
	copy(c.p, path.p)
	return c
}

func (path day16Path) loop(p grid.Pos) bool {
	for _, c := range path.p {
		if c == p {
			return true
		}
	}
	return false
}
