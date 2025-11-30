package main

import (
	"fmt"
	"io"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("20a", day20a)
	registerDay("20b", day20b)
}

func day20a(args []string, rdr io.Reader) (string, error) {
	sum, err := day20Run(rdr, 100, 2)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", sum), nil
}

func day20b(args []string, rdr io.Reader) (string, error) {
	sum, err := day20Run(rdr, 100, 20)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", sum), nil
}

func day20Run(rdr io.Reader, minSave, maxJump int) (int, error) {
	g, err := grid.FromReader(rdr)
	if err != nil {
		return 0, err
	}

	start, end := grid.Pos{}, grid.Pos{}
	for p, r := range g.Walk() {
		if r == 'S' {
			start = p
			g.G[p.Row][p.Col] = '.'
		}
		if r == 'E' {
			end = p
			g.G[p.Row][p.Col] = '.'
		}
	}
	bestStart, err := day20CalcBest(g, start)
	if err != nil {
		return 0, err
	}
	bestEnd, err := day20CalcBest(g, end)
	if err != nil {
		return 0, err
	}
	bestTotal, ok := bestStart[end]
	if !ok {
		return 0, fmt.Errorf("cannot find a path")
	}
	bestGoal := bestTotal - minSave
	sum := 0
	for ps := range g.Walk() {
		bs, ok := bestStart[ps]
		if !ok {
			continue
		}
		for dr := -1 * maxJump; dr <= maxJump; dr++ {
			drAbs := dr
			if drAbs < 0 {
				drAbs *= -1
			}
			for dc := drAbs - maxJump; dc <= maxJump-drAbs; dc++ {
				pe := grid.Pos{Row: ps.Row + dr, Col: ps.Col + dc}
				be, ok := bestEnd[pe]
				if !ok {
					continue
				}
				dyAbs := dc
				if dyAbs < 0 {
					dyAbs *= -1
				}
				if bs+be+drAbs+dyAbs <= bestGoal {
					sum++
				}
			}
		}
	}
	fmt.Printf("grid can be solved in %d, found %d jumps that save %d time in %d jumps\n", bestTotal, sum, minSave, maxJump)
	return sum, nil
}

func day20CalcBest(g grid.Grid, start grid.Pos) (map[grid.Pos]int, error) {
	best := map[grid.Pos]int{start: 0}
	if !g.ValidPos(start) {
		return best, fmt.Errorf("invalid start location")
	}
	search := []grid.Pos{start}
	for len(search) > 0 {
		cur := search[0]
		for d := range grid.DirIter() {
			check := cur.MoveD(d)
			if !g.ValidPos(check) || g.G[check.Row][check.Col] == '#' {
				continue
			}
			if nb, ok := best[check]; !ok || nb > best[cur]+1 {
				best[check] = best[cur] + 1
				search = append(search, check)
			}
		}
		// pop the head
		search = search[1:]
	}
	return best, nil
}
