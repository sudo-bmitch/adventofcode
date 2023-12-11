package main

import (
	"fmt"
	"io"
	"strings"
)

type day11Map struct {
	field      [][]rune
	maxX, maxY int
}

type day11Pos struct {
	x, y int
}

func day11a(args []string, rdr io.Reader) (string, error) {
	m, err := day11Parse(rdr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", day11Calc(m, 2)), nil
}

func day11b(args []string, rdr io.Reader) (string, error) {
	m, err := day11Parse(rdr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", day11Calc(m, 1000000)), nil
}

func day11Parse(rdr io.Reader) (day11Map, error) {
	m := day11Map{}
	input, err := io.ReadAll(rdr)
	if err != nil {
		return m, err
	}
	lines := strings.Split(string(input), "\n")
	m.maxX = len(lines)
	m.maxY = len(strings.TrimSpace(lines[0]))
	m.field = make([][]rune, m.maxX)
	for x, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			m.maxX--
			continue
		}
		if m.maxY == 0 {
			m.maxY = len(line)
		}
		if len(line) != m.maxY {
			return m, fmt.Errorf("field width inconsistent: maxY = %d, cur width = %d, line %d", m.maxY, len(line), x)
		}
		m.field[x] = []rune(line)
	}
	return m, nil
}

func day11SumRange(list []int) int {
	sum := 0
	for _, i := range list {
		sum += i
	}
	return sum
}

func day11Calc(m day11Map, expand int) int {
	// calculate the size of space
	spaceX := make([]int, m.maxX)
	spaceY := make([]int, m.maxY)
	for x := range spaceX {
		spaceX[x] = expand
	}
	for y := range spaceY {
		spaceY[y] = expand
	}
	for x := range spaceX {
		for y := range spaceY {
			if m.field[x][y] != '.' {
				spaceX[x] = 1
				spaceY[y] = 1
			}
		}
	}
	// list all galaxies
	galaxies := []day11Pos{}
	for x := range spaceX {
		for y := range spaceY {
			if m.field[x][y] != '.' {
				galaxies = append(galaxies, day11Pos{x: x, y: y})
			}
		}
	}
	// sum the distances between each
	sum := 0
	for i, g1 := range galaxies {
		if i >= len(galaxies)-1 {
			break
		}
		for _, g2 := range galaxies[i+1:] {
			d := 0
			if g1.x < g2.x {
				d += day11SumRange(spaceX[g1.x:g2.x])
			} else if g1.x > g2.x {
				d += day11SumRange(spaceX[g2.x:g1.x])
			}
			if g1.y < g2.y {
				d += day11SumRange(spaceY[g1.y:g2.y])
			} else if g1.y > g2.y {
				d += day11SumRange(spaceY[g2.y:g1.y])
			}
			sum += d
		}
	}
	return sum
}
