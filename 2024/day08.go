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
	for p, v := range g.Walk() {
		if v == '.' {
			continue
		}
		antennas[v] = append(antennas[v], p)
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
				dr, dc := pos1.Row-pos2.Row, pos1.Col-pos2.Col
				node1 := grid.Pos{Row: pos1.Row + dr, Col: pos1.Col + dc}
				if antinodes.ValidPos(node1) {
					antinodes.G[node1.Row][node1.Col] = '#'
				}
				node2 := grid.Pos{Row: pos2.Row - dr, Col: pos2.Col - dc}
				if antinodes.ValidPos(node2) {
					antinodes.G[node2.Row][node2.Col] = '#'
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
				dr, dc := pos1.Row-pos2.Row, pos1.Col-pos2.Col
				for i := 0; ; i++ {
					node := grid.Pos{Row: pos1.Row + (dr * i), Col: pos1.Col + (dc * i)}
					if !antinodes.ValidPos(node) {
						break
					}
					antinodes.G[node.Row][node.Col] = '#'
				}
				for i := -1; ; i-- {
					node := grid.Pos{Row: pos1.Row + (dr * i), Col: pos1.Col + (dc * i)}
					if !antinodes.ValidPos(node) {
						break
					}
					antinodes.G[node.Row][node.Col] = '#'
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
