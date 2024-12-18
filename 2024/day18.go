package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("18a", day18a)
	registerDay("18b", day18b)
}

func day18a(args []string, rdr io.Reader) (string, error) {
	sum, err := day18Run(71, 71, rdr, 1024)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", sum), nil
}

func day18b(args []string, rdr io.Reader) (string, error) {
	return day18FindLimit(71, 71, rdr)
}

func day18FindLimit(w, h int, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(in), "\n")
	min, max := 0, len(lines)
	for max-min > 1 {
		mid := min + (max-min)/2
		_, err = day18Run(w, h, bytes.NewBuffer(in), mid)
		if err != nil {
			// too long, drop max
			max = mid
		} else {
			// too short, increase min
			min = mid
		}
	}
	return lines[max-1], nil
}

func day18Run(w, h int, rdr io.Reader, limit int) (int, error) {
	// note width/height swapped for grid library
	g, err := grid.New(h, w)
	if err != nil {
		return 0, err
	}
	for p := range g.Walk() {
		g.G[p.X][p.Y] = '.'
	}
	in := bufio.NewScanner(rdr)
	lines := 0
	for in.Scan() {
		// parse each line into a solution and list of numbers
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		nums := strings.Split(line, ",")
		if len(nums) != 2 {
			return 0, fmt.Errorf("invalid line: %s", line)
		}
		x, err := strconv.Atoi(nums[0])
		if err != nil {
			return 0, fmt.Errorf("failed converting x(%s) in line %s: %w", nums[0], line, err)
		}
		y, err := strconv.Atoi(nums[1])
		if err != nil {
			return 0, fmt.Errorf("failed converting y(%s) in line %s: %w", nums[1], line, err)
		}
		if x > h || y > w {
			return 0, fmt.Errorf("coordinates out of range: %d,%d on line %s", x, y, line)
		}
		g.G[x][y] = '#'
		// only read lines up to limit
		lines++
		if lines >= limit {
			break
		}
	}
	// bfs over grid for shortest path from 0,0 to w,h
	goal := grid.Pos{X: h - 1, Y: w - 1}
	paths := []day18Path{
		{{X: 0, Y: 0}},
	}
	best := map[grid.Pos]int{
		{X: 0, Y: 0}: 1,
	}
	for len(paths) > 0 {
		// sort by path length, always try to move the shortest entry first
		slices.SortFunc(paths, func(a, b day18Path) int {
			if len(a) < len(b) {
				return -1
			}
			if len(a) > len(b) {
				return 1
			}
			return 0
		})
		// for the head entry, try each direction
		curPos := paths[0][len(paths[0])-1]
		for d := grid.North; d < grid.DirLen; d++ {
			nextPos := curPos.MoveD(d)
			if !g.ValidPos(nextPos) || g.G[nextPos.X][nextPos.Y] == '#' || !paths[0].valid(nextPos) {
				// skip entries that leave the map, hit a wall, or loop
				continue
			}
			if b, ok := best[nextPos]; ok && b <= len(paths[0])+1 {
				// skip entries that are worse than current best to this position
				continue
			}
			// add a path
			newPath := paths[0].clone()
			newPath = append(newPath, nextPos)
			// check for goal
			if nextPos == goal {
				// for _, p := range newPath {
				// 	g.G[p.X][p.Y] = 'O'
				// }
				// fmt.Printf("found a path in %d steps:\n%s\n", len(newPath)-1, g.String())
				return len(newPath) - 1, nil
			}
			paths = append(paths, newPath)
			// track current best to this position
			best[nextPos] = len(newPath)
		}
		// drop the head entry
		paths = paths[1:]

	}
	return 0, fmt.Errorf("no paths found")
}

type day18Path []grid.Pos

func (path day18Path) clone() day18Path {
	ret := make([]grid.Pos, len(path))
	copy(ret, path)
	return ret
}

func (path day18Path) valid(p grid.Pos) bool {
	for _, e := range path {
		if p == e {
			return false
		}
	}
	return true
}
