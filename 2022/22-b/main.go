package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

var debug = true

type pos int

const (
	posEmpty pos = iota
	posWalk
	posSolid
)

type dir int

const (
	dirRight dir = iota
	dirDown
	dirLeft
	dirUp
	dirCount
)

type edge struct {
	side *cubeSide // link to the other side, to get the posD/posR
	dir  dir       // which of the 4 directions does this connect: r, d, l, r
}
type cubeSide struct {
	edges      [4]edge // r, d, l, u
	posD, posR int     // top left corner of the cube
}

type board struct {
	g       [][]pos
	w, h    int
	sides   [][]*cubeSide
	edgeLen int
}

type cursor struct {
	posD, posR int // position
	d          dir // direction facing
}

type move struct {
	distance int
	turn     int // L = -1, R = 1, none = 0
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	lineBoardRE := regexp.MustCompile(`^[ \.#]+$`)
	linePathRE := regexp.MustCompile(`^[0-9LR]+$`)
	b := board{g: [][]pos{}}
	path := []move{}
	// parse input
	for in.Scan() {
		// line := strings.TrimSpace(in.Text()) // leading whitespace is needed
		line := in.Text()
		if line == "" {
			continue
		}
		matchBoard := lineBoardRE.Match([]byte(line))
		matchPath := linePathRE.Match([]byte(line))
		if matchBoard && len(path) == 0 {
			row := []pos{}
			for _, c := range line {
				switch c {
				case ' ':
					row = append(row, posEmpty)
				case '.':
					row = append(row, posWalk)
				case '#':
					row = append(row, posSolid)
				default:
					fmt.Fprintf(os.Stderr, "failed to parse %s\n", line)
					return
				}
			}
			b.g = append(b.g, row)
			if len(row) > b.w {
				b.w = len(row)
			}
		} else if matchPath {
			distance := 0
			for _, c := range line {
				switch {
				case c >= rune('0') && c <= rune('9'):
					distance = distance*10 + int(c-rune('0'))
				case c == 'L':
					path = append(path, move{distance: distance, turn: -1})
					distance = 0
				case c == 'R':
					path = append(path, move{distance: distance, turn: 1})
					distance = 0
				default:
					fmt.Fprintf(os.Stderr, "failed to parse %s\n", line)
					return
				}
			}
			if distance > 0 {
				path = append(path, move{distance: distance})
			}
		} else {
			fmt.Fprintf(os.Stderr, "failed to parse %s\n", line)
			return
		}
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	b.h = len(b.g)
	if b.h == 0 {
		fmt.Fprintf(os.Stderr, "failed to read any board\n")
		return
	}

	// define cube shape, hard coded because monkeys are impatient
	if b.h == 12 && b.w == 16 {
		// xxAx
		// BCDx
		// xxEF
		sideA := &cubeSide{posD: 0, posR: 2}
		sideB := &cubeSide{posD: 1, posR: 0}
		sideC := &cubeSide{posD: 1, posR: 1}
		sideD := &cubeSide{posD: 1, posR: 2}
		sideE := &cubeSide{posD: 2, posR: 2}
		sideF := &cubeSide{posD: 2, posR: 3}
		sideA.edges = [4]edge{{side: sideF, dir: dirLeft}, {side: sideD, dir: dirDown}, {side: sideC, dir: dirDown}, {side: sideB, dir: dirRight}}
		sideB.edges = [4]edge{{side: sideC, dir: dirRight}, {side: sideE, dir: dirUp}, {side: sideF, dir: dirUp}, {side: sideA, dir: dirDown}}
		sideC.edges = [4]edge{{side: sideD, dir: dirRight}, {side: sideE, dir: dirRight}, {side: sideB, dir: dirLeft}, {side: sideA, dir: dirRight}}
		sideD.edges = [4]edge{{side: sideF, dir: dirDown}, {side: sideE, dir: dirDown}, {side: sideC, dir: dirLeft}, {side: sideA, dir: dirUp}}
		sideE.edges = [4]edge{{side: sideF, dir: dirRight}, {side: sideB, dir: dirUp}, {side: sideC, dir: dirUp}, {side: sideD, dir: dirUp}}
		sideF.edges = [4]edge{{side: sideA, dir: dirRight}, {side: sideB, dir: dirRight}, {side: sideE, dir: dirLeft}, {side: sideD, dir: dirLeft}}
		b.edgeLen = 4
		b.sides = [][]*cubeSide{
			{nil, nil, sideA, nil},
			{sideB, sideC, sideD, nil},
			{nil, nil, sideE, sideF},
		}
	} else if b.h == 200 && b.w == 150 {
		// xAB
		// xCx
		// DEx
		// Fxx
		sideA := &cubeSide{posD: 0, posR: 1}
		sideB := &cubeSide{posD: 0, posR: 2}
		sideC := &cubeSide{posD: 1, posR: 1}
		sideD := &cubeSide{posD: 2, posR: 0}
		sideE := &cubeSide{posD: 2, posR: 1}
		sideF := &cubeSide{posD: 3, posR: 0}
		sideA.edges = [4]edge{{side: sideB, dir: dirRight}, {side: sideC, dir: dirDown}, {side: sideD, dir: dirRight}, {side: sideF, dir: dirRight}}
		sideB.edges = [4]edge{{side: sideE, dir: dirLeft}, {side: sideC, dir: dirLeft}, {side: sideA, dir: dirLeft}, {side: sideF, dir: dirUp}}
		sideC.edges = [4]edge{{side: sideB, dir: dirUp}, {side: sideE, dir: dirDown}, {side: sideD, dir: dirDown}, {side: sideA, dir: dirUp}}
		sideD.edges = [4]edge{{side: sideE, dir: dirRight}, {side: sideF, dir: dirDown}, {side: sideA, dir: dirRight}, {side: sideC, dir: dirRight}}
		sideE.edges = [4]edge{{side: sideB, dir: dirLeft}, {side: sideF, dir: dirLeft}, {side: sideD, dir: dirLeft}, {side: sideC, dir: dirUp}}
		sideF.edges = [4]edge{{side: sideE, dir: dirUp}, {side: sideB, dir: dirDown}, {side: sideA, dir: dirDown}, {side: sideD, dir: dirUp}}
		b.edgeLen = 50
		b.sides = [][]*cubeSide{
			{nil, sideA, sideB},
			{nil, sideC, nil},
			{sideD, sideE, nil},
			{sideF, nil, nil},
		}

	} else {
		fmt.Fprintf(os.Stderr, "unknown size to map to cube: %d x %d\n", b.h, b.w)
		return
	}

	// get initial position
	c := cursor{
		posD: 0,
		d:    dirRight,
	}
	for i := range b.g[0] {
		if b.g[0][i] == posWalk {
			c.posR = i
			break
		}
	}
	for _, p := range path {
		b.walk(&c, p)
		if debug {
			fmt.Printf("moved %d and turned %d, new pos %d,%d facing %d\n", p.distance, p.turn, c.posD, c.posR, c.d)
		}
	}

	result := 1000*(c.posD+1) + 4*(c.posR+1) + int(c.d)
	fmt.Printf("Result: %d\n", result)
}

func (b board) walk(c *cursor, m move) {
	for i := 0; i < m.distance; i++ {
		blocked, pD, pR, nD := b.next(*c)
		if blocked {
			break
		}
		c.posD, c.posR, c.d = pD, pR, nD
	}
	if m.turn != 0 {
		c.d = dir((int(c.d) + m.turn + int(dirCount)) % int(dirCount))
	}
}

// next returns if next step is blocked, and the resulting D/R/dir
func (b board) next(c cursor) (bool, int, int, dir) {
	dD, dR := 0, 0
	switch c.d {
	case dirRight:
		dR = 1
	case dirDown:
		dD = 1
	case dirLeft:
		dR = -1
	case dirUp:
		dD = -1
	default:
		fmt.Fprintf(os.Stderr, "unknown direction %d\n", c.d)
		os.Exit(1)
	}
	nextD, nextR, nextDir := c.posD+dD, c.posR+dR, c.d
	if nextD < 0 || nextR < 0 || nextD >= len(b.g) || nextR >= len(b.g[nextD]) || b.g[nextD][nextR] == posEmpty {
		// need to handle the cube fold
		if debug {
			fmt.Printf("reached edge at %dx%d facing %d\n", nextD, nextR, nextDir)
		}
		// identify my current side
		pD := c.posD / b.edgeLen
		pR := c.posR / b.edgeLen
		// get offset from my left side of where I'm facing
		offset := 0
		switch c.d {
		case dirRight:
			offset = c.posD % b.edgeLen
		case dirDown:
			offset = b.edgeLen - 1 - c.posR%b.edgeLen
		case dirLeft:
			offset = b.edgeLen - 1 - c.posD%b.edgeLen
		case dirUp:
			offset = c.posR % b.edgeLen
		}
		// follow the edge
		curSide := b.sides[pD][pR]
		curEdge := curSide.edges[c.d]
		nextDir = curEdge.dir
		switch nextDir {
		case dirRight:
			nextR = curEdge.side.posR * b.edgeLen
			nextD = curEdge.side.posD*b.edgeLen + offset
		case dirDown:
			nextR = curEdge.side.posR*b.edgeLen + b.edgeLen - 1 - offset
			nextD = curEdge.side.posD * b.edgeLen
		case dirLeft:
			nextR = curEdge.side.posR*b.edgeLen + b.edgeLen - 1
			nextD = curEdge.side.posD*b.edgeLen + b.edgeLen - 1 - offset
		case dirUp:
			nextR = curEdge.side.posR*b.edgeLen + offset
			nextD = curEdge.side.posD*b.edgeLen + b.edgeLen - 1
		}
		if debug {
			fmt.Printf("wrapped edge to %dx%d facing %d\n", nextD, nextR, nextDir)
		}
	}
	if b.g[nextD][nextR] == posSolid {
		return true, c.posD, c.posR, c.d
	} else {
		return false, nextD, nextR, nextDir
	}
}

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// func mustAtoi(s string) int {
// 	i, err := strconv.Atoi(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return i
// }
