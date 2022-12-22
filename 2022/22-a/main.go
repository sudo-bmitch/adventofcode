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

type board struct {
	g [][]pos
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
		// line := strings.TrimSpace(in.Text()) // leading whitespace is needed, but the puzzle apparently worked with this, I don't know how
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

	if len(b.g) == 0 {
		fmt.Fprintf(os.Stderr, "failed to read any board\n")
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
		blocked, pD, pR := b.next(*c)
		if blocked {
			break
		}
		c.posD, c.posR = pD, pR
	}
	if m.turn != 0 {
		c.d = dir((int(c.d) + m.turn + int(dirCount)) % int(dirCount))
	}
}

// next returns if next step is blocked, and the resulting D/R position in that dir
func (b board) next(c cursor) (bool, int, int) {
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
	nextD, nextR := c.posD, c.posR
	for {
		nextD = (nextD + dD + len(b.g)) % len(b.g)
		nextR = (nextR + dR + len(b.g[nextD])) % len(b.g[nextD])
		if b.g[nextD][nextR] == posSolid || nextD == c.posD && nextR == c.posR {
			// if blocked or looped around entire map
			return true, c.posD, c.posR
		} else if b.g[nextD][nextR] == posWalk {
			// found next walkable position
			return false, nextD, nextR
		}
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
