// Package grid is used for 2d grids.
package grid

import (
	"fmt"
	"io"
	"strings"
)

type Grid struct {
	G    [][]rune
	H, W int
}

func FromReader(rdr io.Reader) (Grid, error) {
	b, err := io.ReadAll(rdr)
	if err != nil {
		return Grid{}, err
	}
	return FromString(string(b))
}

func FromString(s string) (Grid, error) {
	rows := strings.Split(s, "\n")
	g := Grid{
		G: [][]rune{},
	}
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if row == "" {
			continue
		}
		if g.W == 0 {
			g.W = len(row)
		} else if g.W != len(row) {
			return g, fmt.Errorf("inconsistent row width, %d != %d, line = %s", g.W, len(row), row)
		}
		g.G = append(g.G, []rune(row))
	}
	g.H = len(g.G)
	return g, nil
}

func (g Grid) ValidPos(p Pos) bool {
	if p.X >= 0 && p.Y >= 0 && p.X < g.H && p.Y < g.W {
		return true
	}
	return false
}

type Pos struct {
	X, Y int
}

func (p Pos) MoveP(m Pos) Pos {
	return Pos{X: p.X + m.X, Y: p.Y + m.Y}
}

func (p Pos) MoveD(d Dir) Pos {
	return p.MoveP(DirPos[d])
}

func (p Pos) String() string {
	return fmt.Sprintf("[%d,%d]", p.X, p.Y)
}

type Dir int

const (
	North Dir = iota
	East
	South
	West
	DirLen
)

var DirPos = [DirLen]Pos{
	North: {X: -1},
	East:  {Y: 1},
	South: {X: 1},
	West:  {Y: -1},
}
