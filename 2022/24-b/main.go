package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

var debug = false

type pos struct {
	x, y int
}

type storm struct {
	pos pos // absolute
	dir pos // delta
}

type space int

const (
	spaceOpen space = iota
	spaceWall
	spaceStorm
)

type board struct {
	grid   [][]space
	storms []storm
}

func main() {
	b := board{grid: [][]space{}, storms: []storm{}}
	in := bufio.NewScanner(os.Stdin)
	lineRE := regexp.MustCompile(`^[\.<^>v#]+$`)
	// parse input
	for in.Scan() {
		line := in.Text()
		if line == "" {
			continue
		}
		if !lineRE.MatchString(line) {
			fmt.Fprintf(os.Stderr, "failed to parse %s\n", line)
			return
		}
		if len(b.grid) > 0 && len(line) != len(b.grid[0]) {
			fmt.Fprintf(os.Stderr, "line length irregular on %s\n", line)
			return
		}
		b.grid = append(b.grid, make([]space, len(line)))
		x := len(b.grid) - 1
		for y, c := range line {
			curPos := pos{x: x, y: y}
			switch c {
			case '#':
				b.grid[x][y] = spaceWall
			case '>':
				b.grid[x][y] = spaceStorm
				b.storms = append(b.storms, storm{pos: curPos, dir: pos{y: 1}})
			case '<':
				b.grid[x][y] = spaceStorm
				b.storms = append(b.storms, storm{pos: curPos, dir: pos{y: -1}})
			case '^':
				b.grid[x][y] = spaceStorm
				b.storms = append(b.storms, storm{pos: curPos, dir: pos{x: -1}})
			case 'v':
				b.grid[x][y] = spaceStorm
				b.storms = append(b.storms, storm{pos: curPos, dir: pos{x: 1}})
			}
		}
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	elvesCur := map[pos]bool{{x: 0, y: 1}: true} // start at 0,1
	goals := [3]pos{{x: len(b.grid) - 1, y: len(b.grid[0]) - 2}, {x: 0, y: 1}, {x: len(b.grid) - 1, y: len(b.grid[0]) - 2}}
	result := 0

	for _, goal := range goals {
		goalReached := false
		t := 0
		for t = 0; !goalReached; t++ {
			elvesNext := map[pos]bool{}
			// move the storms
			b = b.next()
			// add possible positions to elvesNext
			for p := range elvesCur {
				// try each direction
				for _, delta := range [5]pos{{}, {x: -1}, {x: 1}, {y: -1}, {y: 1}} {
					pd := p
					pd.x += delta.x
					pd.y += delta.y
					// fail if space leaves grid, is a wall, or a storm
					if pd.x < 0 || pd.y < 0 || pd.x >= len(b.grid) || pd.y >= len(b.grid[0]) || b.grid[pd.x][pd.y] != spaceOpen {
						continue
					}
					elvesNext[pd] = true
					// stop if any elf makes it to goal
					if pd.x == goal.x && pd.y == goal.y {
						goalReached = true
					}
				}
			}
			elvesCur = elvesNext
			if debug {
				fmt.Printf("board at time %d\n", t)
				b.print(elvesCur)
			}
		}
		elvesCur = map[pos]bool{goal: true} // reset to start from new goal
		result += t
	}

	fmt.Printf("Result: %d\n", result)
}

func (b board) next() board {
	for x := 0; x < len(b.grid); x++ {
		for y := 0; y < len(b.grid[x]); y++ {
			if b.grid[x][y] == spaceStorm {
				b.grid[x][y] = spaceOpen // reset all the spaces
			}
		}
	}
	// move each of the storms
	for i := range b.storms {
		b.storms[i].pos.x += b.storms[i].dir.x
		b.storms[i].pos.y += b.storms[i].dir.y
		if b.storms[i].pos.x < 1 {
			b.storms[i].pos.x = len(b.grid) - 2
		}
		if b.storms[i].pos.x > len(b.grid)-2 {
			b.storms[i].pos.x = 1
		}
		if b.storms[i].pos.y < 1 {
			b.storms[i].pos.y = len(b.grid[0]) - 2
		}
		if b.storms[i].pos.y > len(b.grid[0])-2 {
			b.storms[i].pos.y = 1
		}
		b.grid[b.storms[i].pos.x][b.storms[i].pos.y] = spaceStorm
	}
	return b
}

func (b board) print(e map[pos]bool) {
	for x := 0; x < len(b.grid); x++ {
		for y := 0; y < len(b.grid[x]); y++ {
			switch b.grid[x][y] {
			case spaceOpen:
				if e[pos{x: x, y: y}] {
					fmt.Printf("e")
				} else {
					fmt.Printf(".")
				}
			case spaceWall:
				fmt.Printf("#")
			case spaceStorm:
				fmt.Printf("@")
			default:
				fmt.Printf("?")
			}
		}
		fmt.Printf("\n")
	}
}

// func min(a, b int) int {
// 	if a > b {
// 		return b
// 	}
// 	return a
// }
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
