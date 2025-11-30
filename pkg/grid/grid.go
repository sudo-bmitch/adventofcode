// Package grid is used for 2d grids.
package grid

import (
	"fmt"
	"io"
	"iter"
	"strings"
)

type Grid struct {
	G    [][]rune
	H, W int
}

func New(w, h int) (Grid, error) {
	if w <= 0 || h <= 0 {
		return Grid{}, fmt.Errorf("grid width and height must be positive")
	}
	g := Grid{
		G: make([][]rune, h),
		W: w,
		H: h,
	}
	for i := range g.G {
		g.G[i] = make([]rune, w)
	}
	return g, nil
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

func (g Grid) String() string {
	buf := strings.Builder{}
	for _, row := range g.G {
		buf.WriteString(string(row))
		buf.WriteRune('\n')
	}
	return buf.String()
}

func (g Grid) ValidPos(p Pos) bool {
	if p.Row >= 0 && p.Col >= 0 && p.Row < g.H && p.Col < g.W {
		return true
	}
	return false
}

func (g Grid) Walk() iter.Seq2[Pos, rune] {
	return func(yield func(Pos, rune) bool) {
		for r := 0; r < g.H; r++ {
			for c := 0; c < g.W; c++ {
				if !yield(Pos{Row: r, Col: c}, g.G[r][c]) {
					return
				}
			}
		}
	}
}

type Pos struct {
	Row int
	Col int
}

func (p Pos) MoveP(m Pos) Pos {
	return Pos{Row: p.Row + m.Row, Col: p.Col + m.Col}
}

func (p Pos) MoveD(d Dir) Pos {
	return p.MoveP(DirPos[d])
}

func (p Pos) String() string {
	return fmt.Sprintf("[%d,%d]", p.Row, p.Col)
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
	North: {Row: -1},
	East:  {Col: 1},
	South: {Row: 1},
	West:  {Col: -1},
}

func (d Dir) String() string {
	switch d {
	case North:
		return "north"
	case East:
		return "east"
	case South:
		return "south"
	case West:
		return "west"
	default:
		return "unknown"
	}
}

func DirIter() iter.Seq[Dir] {
	return func(yield func(Dir) bool) {
		for d := range DirLen {
			if !yield(d) {
				return
			}
		}
	}
}
