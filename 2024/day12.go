package main

import (
	"fmt"
	"io"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("12a", day12a)
	registerDay("12b", day12b)
}

func day12a(args []string, rdr io.Reader) (string, error) {
	regions, seen, err := day12Regions(rdr)
	if err != nil {
		return "", err
	}

	// compute the area * perimeter of each region
	sum := 0
	for _, region := range regions {
		regionNum := seen.G[region[0].Row][region[0].Col]
		perimeter := 0
		area := len(region)
		for _, cur := range region {
			for d := range grid.DirIter() {
				try := cur.MoveD(d)
				if !seen.ValidPos(try) || seen.G[try.Row][try.Col] != regionNum {
					perimeter++
				}
			}
		}
		sum += (area * perimeter)
	}
	return fmt.Sprintf("%d", sum), nil
}

func day12b(args []string, rdr io.Reader) (string, error) {
	regions, seen, err := day12Regions(rdr)
	if err != nil {
		return "", err
	}

	// compute the area * sides of each region
	sum := 0
	// map is used to exclude sides that will be counted by another position
	sideMap := map[grid.Dir]grid.Dir{
		grid.North: grid.West,
		grid.East:  grid.North,
		grid.South: grid.East,
		grid.West:  grid.South,
	}
	for _, region := range regions {
		regionNum := seen.G[region[0].Row][region[0].Col]
		sides := 0
		area := len(region)
		for _, cur := range region {
			for d := range grid.DirIter() {
				try := cur.MoveD(d)
				prevSideA := cur.MoveD(sideMap[d])
				prevSideB := try.MoveD(sideMap[d])
				if (!seen.ValidPos(try) || seen.G[try.Row][try.Col] != regionNum) && // fence needed
					((!seen.ValidPos(prevSideA) || seen.G[prevSideA.Row][prevSideA.Col] != regionNum) || // previous does not match
						(seen.ValidPos(prevSideA) && seen.G[prevSideA.Row][prevSideA.Col] == regionNum && seen.ValidPos(prevSideB) && seen.G[prevSideB.Row][prevSideB.Col] == regionNum)) { // inside corner
					sides++
				}
			}
		}
		sum += (area * sides)
	}
	return fmt.Sprintf("%d", sum), nil
}

func day12Regions(rdr io.Reader) ([][]grid.Pos, grid.Grid, error) {
	g, err := grid.FromReader(rdr)
	if err != nil {
		return nil, grid.Grid{}, err
	}

	// break up into regions
	regions := [][]grid.Pos{}
	seen, err := grid.New(g.W, g.H)
	if err != nil {
		return nil, grid.Grid{}, err
	}
	for p, r := range g.Walk() {
		if seen.G[p.Row][p.Col] != 0 {
			continue
		}
		queue := []grid.Pos{p}
		region := []grid.Pos{}
		regionNum := rune(len(regions) + 1)
		seen.G[p.Row][p.Col] = regionNum
		for len(queue) > 0 {
			cur := queue[len(queue)-1]
			queue = queue[:len(queue)-1]
			region = append(region, cur)
			// add neighbors to queue that haven't been seen yet
			for d := range grid.DirIter() {
				try := cur.MoveD(d)
				if g.ValidPos(try) && g.G[try.Row][try.Col] == r && seen.G[try.Row][try.Col] == 0 {
					queue = append(queue, try)
					seen.G[try.Row][try.Col] = regionNum
				}
			}
		}
		regions = append(regions, region)
	}
	return regions, seen, nil
}
