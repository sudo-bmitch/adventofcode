package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x, y int
}

func newPoint(str string) (point, error) {
	posStr := strings.SplitN(str, ",", 2)
	if len(posStr) < 2 {
		return point{}, fmt.Errorf("comma not found")
	}
	y, err := strconv.Atoi(posStr[0])
	if err != nil {
		return point{}, fmt.Errorf("failed to parse y: %w", err)
	}
	x, err := strconv.Atoi(posStr[1])
	if err != nil {
		return point{}, fmt.Errorf("failed to parse x: %w", err)
	}
	return point{x: x, y: y}, nil
}

var sandStart = point{y: 500, x: 0}

type content int

const (
	void content = iota
	air
	rock
	sand
)

type grid struct {
	g                      [][]content
	xMax, xMin, yMax, yMin int
}

func (g *grid) addRock(a, b point) error {
	// validate a and b, no negative, no diagonal
	if a.x < 0 || a.y < 0 || b.x < 0 || b.y < 0 {
		return fmt.Errorf("negative point values not allowed")
	}
	if a.x != b.x && a.y != b.y {
		return fmt.Errorf("diagonal not allowed")
	}
	xMax := max(a.x, b.x)
	xMin := min(a.x, b.x)
	yMax := max(a.y, b.y)
	yMin := min(a.y, b.y)
	if g.xMax == 0 && g.xMin == 0 && g.yMax == 0 && g.yMin == 0 {
		// initialize
		g.xMax = xMax
		g.xMin = xMin
		g.yMax = yMax
		g.yMin = yMin
	}
	g.xMax = max(g.xMax, xMax)
	g.xMin = min(g.xMin, xMin)
	g.yMax = max(g.yMax, yMax)
	g.yMin = min(g.yMin, yMin)
	// if max x/y > grid size, extend
	g.expand(point{x: xMax, y: yMax})
	// identify if this is horizontal or vertical, get the min and max of the other axis, set positions to rock
	if a.x == b.x {
		// horizontal
		for y := yMin; y <= yMax; y++ {
			g.g[a.x][y] = rock
			g.addAir(point{x: a.x, y: y})
		}
	} else {
		// vertical
		for x := xMin; x <= xMax; x++ {
			g.g[x][a.y] = rock
			g.addAir(point{x: x, y: a.y})
		}
	}
	return nil
}

func (g *grid) addAir(p point) {
	for x := p.x - 1; x >= 0; x-- {
		if g.g[x][p.y] != void {
			break
		}
		g.g[x][p.y] = air
	}
}

func (g *grid) dropSand(p point) error {
	// if not starting in air, fail
	if p.x < 0 || p.y < 0 || p.x >= len(g.g) || p.y >= len(g.g[0]) || g.g[p.x][p.y] != air {
		return fmt.Errorf("not starting in air")
	}
	for {
		// reached bottom or side of grid, only void below
		if p.x >= len(g.g)-1 || p.y <= 0 || p.y >= len(g.g[0]) {
			return fmt.Errorf("off bottom/side of grid")
		}
		// try down
		if g.g[p.x+1][p.y] == air {
			p.x++
			continue
		}
		if g.g[p.x+1][p.y] == void {
			return fmt.Errorf("fell into void")
		}
		// try down left
		if g.g[p.x+1][p.y-1] == air {
			p.x++
			p.y--
			continue
		}
		if g.g[p.x+1][p.y-1] == void {
			return fmt.Errorf("fell into void")
		}
		// try down right
		if g.g[p.x+1][p.y+1] == air {
			p.x++
			p.y++
			continue
		}
		if g.g[p.x+1][p.y+1] == void {
			return fmt.Errorf("fell into void")
		}
		// reached bottom, place sand
		g.g[p.x][p.y] = sand
		return nil
	}
}

func (g *grid) expand(p point) {
	yLen := p.y + 1
	if len(g.g) > 0 && yLen < len(g.g[0]) {
		yLen = len(g.g[0])
	}
	for x, row := range g.g {
		if len(row) < yLen {
			// extend columns
			add := make([]content, yLen-len(row))
			g.g[x] = append(g.g[x], add...)
		}
	}
	for len(g.g) < (p.x + 1) {
		// append rows
		g.g = append((g.g), make([]content, yLen))
	}
}

func (g *grid) print() {
	chars := []rune{
		'~', '.', '#', 'o',
	}
	for _, row := range g.g {
		for y := g.yMin; y <= g.yMax; y++ {
			c := row[y]
			fmt.Printf("%c", chars[c])
		}
		fmt.Printf("\n")
	}
}

func main() {
	g := grid{}
	g.expand(sandStart)
	in := bufio.NewScanner(os.Stdin)
	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		pointStrList := strings.Split(line, " -> ")
		prevPoint := point{}
		for i, pointStr := range pointStrList {
			curPoint, err := newPoint(pointStr)
			if err != nil {
				panic(err)
			}
			if i > 0 {
				err := g.addRock(prevPoint, curPoint)
				if err != nil {
					panic(err)
				}
			}
			prevPoint = curPoint
		}
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}
	// add floor
	err := g.addRock(point{x: g.xMax + 2, y: 500 - g.xMax - 2}, point{x: g.xMax + 2, y: 500 + g.xMax + 2})
	if err != nil {
		fmt.Printf("failed to add floor: %v\n", err)
		return
	}
	// add sand
	count := 0
	for {
		if err := g.dropSand(sandStart); err != nil {
			break
		}
		count++
	}
	// show grid
	fmt.Printf("End with sand:\n")
	g.print()
	fmt.Printf("Result: %d\n", count)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
