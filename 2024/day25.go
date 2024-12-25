package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("25a", day25a)
	registerDay("25b", day25b)
}

func day25a(args []string, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	schematics := strings.Split(string(in), "\n\n")
	keys := [][]int{}
	locks := [][]int{}
	for _, schematic := range schematics {
		schematic = strings.TrimSpace(schematic)
		if schematic == "" {
			continue
		}
		g, err := grid.FromString(schematic)
		if err != nil {
			return "", err
		}
		if schematic[:5] == "#####" {
			// reading a lock
			lock := make([]int, 5)
			for y := 0; y < g.W; y++ {
				for x := g.H - 1; x > 0; x-- {
					if g.G[x][y] == '#' {
						lock[y] = x
						break
					}
				}
			}
			locks = append(locks, lock)
		} else if schematic[:5] == "....." {
			// reading a key
			key := make([]int, 5)
			for y := 0; y < g.W; y++ {
				for x := 0; x < g.H; x++ {
					if g.G[x][y] == '#' {
						key[y] = x
						break
					}
				}
			}
			keys = append(keys, key)
		} else {
			return "", fmt.Errorf("unknown schematic:\n%s", schematic)
		}
	}

	sum := 0
	for _, lock := range locks {
		for _, key := range keys {
			match := true
			for pin := range 5 {
				if lock[pin] >= key[pin] {
					match = false
					break
				}
			}
			if match {
				sum++
			}
		}
	}

	return fmt.Sprintf("%d", sum), nil
}

func day25b(args []string, rdr io.Reader) (string, error) {
	sum := 42

	return fmt.Sprintf("%d", sum), nil
}
