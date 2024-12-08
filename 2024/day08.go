package main

import (
	"fmt"
	"io"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("08a", day08a)
	registerDay("08b", day08b)
}

func day08a(args []string, rdr io.Reader) (string, error) {
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}
	antennas := map[rune][]grid.Pos{}
	for p, r := range g.Walk() {
		if r == '.' {
			continue
		}
		antennas[r] = append(antennas[r], p)
	}
	antinodes, err := grid.New(g.W, g.H)
	if err != nil {
		return "", err
	}
	for freq := range antennas {
		for i, pos1 := range antennas[freq] {
			if i >= len(antennas[freq])-1 {
				continue
			}
			for _, pos2 := range antennas[freq][i+1:] {
				dx, dy := pos1.X-pos2.X, pos1.Y-pos2.Y
				node1 := grid.Pos{X: pos1.X + dx, Y: pos1.Y + dy}
				if antinodes.ValidPos(node1) {
					antinodes.G[node1.X][node1.Y] = '#'
				}
				node2 := grid.Pos{X: pos2.X - dx, Y: pos2.Y - dy}
				if antinodes.ValidPos(node2) {
					antinodes.G[node2.X][node2.Y] = '#'
				}
			}
		}
	}
	sum := 0
	for _, r := range antinodes.Walk() {
		if r == '#' {
			sum++
		}
	}
	return fmt.Sprintf("%d", sum), nil
}
func day08b(args []string, rdr io.Reader) (string, error) {
	g, err := grid.FromReader(rdr)
	if err != nil {
		return "", err
	}
	antennas := map[rune][]grid.Pos{}
	for p, r := range g.Walk() {
		if r == '.' {
			continue
		}
		antennas[r] = append(antennas[r], p)
	}
	antinodes, err := grid.New(g.W, g.H)
	if err != nil {
		return "", err
	}
	for freq := range antennas {
		for i, pos1 := range antennas[freq] {
			if i >= len(antennas[freq])-1 {
				continue
			}
			for _, pos2 := range antennas[freq][i+1:] {
				dx, dy := pos1.X-pos2.X, pos1.Y-pos2.Y
				for i := 0; ; i++ {
					node := grid.Pos{X: pos1.X + (dx * i), Y: pos1.Y + (dy * i)}
					if !antinodes.ValidPos(node) {
						break
					}
					antinodes.G[node.X][node.Y] = '#'
				}
				for i := -1; ; i-- {
					node := grid.Pos{X: pos1.X + (dx * i), Y: pos1.Y + (dy * i)}
					if !antinodes.ValidPos(node) {
						break
					}
					antinodes.G[node.X][node.Y] = '#'
				}
			}
		}
	}
	sum := 0
	for _, r := range antinodes.Walk() {
		if r == '#' {
			sum++
		}
	}
	return fmt.Sprintf("%d", sum), nil
}
