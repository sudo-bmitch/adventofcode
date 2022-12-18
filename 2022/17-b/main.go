package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type move int

const (
	moveLeft move = iota
	moveRight
)

const (
	debug   = false
	dropMax = 1000000000000
	// dropMax      = 2022
	shapeMax     = 4
	chamberWidth = 7
)

type shape struct {
	grid   [shapeMax][shapeMax]bool
	height int
	width  int
}

var shapes = []shape{
	{ // shape: -
		grid: [shapeMax][shapeMax]bool{
			{false, false, false, false},
			{false, false, false, false},
			{false, false, false, false},
			{true, true, true, true},
		},
		height: 1,
		width:  4,
	},
	{ // shape: +
		grid: [shapeMax][shapeMax]bool{
			{false, false, false, false},
			{false, true, false, false},
			{true, true, true, false},
			{false, true, false, false},
		},
		height: 3,
		width:  3,
	},
	{ // shape: L (reverse)
		grid: [shapeMax][shapeMax]bool{
			{false, false, false, false},
			{false, false, true, false},
			{false, false, true, false},
			{true, true, true, false},
		},
		height: 3,
		width:  3,
	},
	{ // shape: |
		grid: [shapeMax][shapeMax]bool{
			{true, false, false, false},
			{true, false, false, false},
			{true, false, false, false},
			{true, false, false, false},
		},
		height: 4,
		width:  1,
	},
	{ // shape: box
		grid: [shapeMax][shapeMax]bool{
			{false, false, false, false},
			{false, false, false, false},
			{true, true, false, false},
			{true, true, false, false},
		},
		height: 2,
		width:  2,
	},
}

type chamber struct {
	rows     [][chamberWidth]bool
	lastDrop []uint64
	lastMove []int
}

func main() {
	c := chamber{
		rows: [][chamberWidth]bool{},
	}
	moves := []move{}
	in := bufio.NewScanner(os.Stdin)
	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		for _, c := range line {
			switch c {
			case '<':
				moves = append(moves, moveLeft)
			case '>':
				moves = append(moves, moveRight)
			default:
				fmt.Fprintf(os.Stderr, "unknown move: %c\n", c)
			}
		}
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	shapeCount := uint64(len(shapes))
	moveCount := len(moves)
	moveI := 0
	totalPattern := uint64(0)
	for drop := uint64(0); drop < dropMax; drop++ {
		s := shapes[drop%shapeCount]
		posU := c.top() + 3
		posR := 2
		for {
			diffR := 0
			switch moves[moveI] {
			case moveLeft:
				diffR = -1
			case moveRight:
				diffR = 1
			}
			moveI = (moveI + 1) % moveCount
			if c.shapeAllowed(s, posU, posR+diffR) {
				posR += diffR
			}
			if !c.shapeAllowed(s, posU-1, posR) {
				break
			}
			posU--
		}
		c.shapePlace(s, posU, posR, moveI, drop)
		if pattern, dropDiff := c.patternSearch(); dropDiff > 0 && dropDiff < dropMax-drop {
			// pattern detected, increment drop and skipped rows counters
			mult := (dropMax - drop) / dropDiff
			drop += (dropDiff * mult)
			totalPattern += (uint64(pattern) * mult)
			fmt.Printf("Skipping %d drops with a %d row pattern repeated %d times and %d drops each\n", (dropDiff * mult), pattern, mult, dropDiff)
		}
		// if drop%1000000 == 0 {
		// 	fmt.Printf("Drop %d, top=%d, cleared=%d\n", drop, c.top(), c.clearedRows)
		// }
	}

	if debug && c.top() < 100 {
		c.print()
	}

	fmt.Printf("Result: %d\n", uint64(c.top())+totalPattern)
}

func (c *chamber) addRows(count int) {
	for i := 0; i < count; i++ {
		c.rows = append(c.rows, [chamberWidth]bool{})
		c.lastMove = append(c.lastMove, 0)
		c.lastDrop = append(c.lastDrop, 0)
	}
}

func (c *chamber) shapeAllowed(s shape, posU int, posR int) bool {
	// calculate if shape is allowed in position
	if posR < 0 || posR+s.width > chamberWidth || posU < 0 {
		return false
	}
	// check for collision with chamber contents
	for u := posU; u < len(c.rows) && u-posU < s.height; u++ {
		// if u >= len(c.rows) {
		// 	continue // no collisions above the top
		// }
		su := shapeMax - (u - posU) - 1 // shape row is inverted (top/bottom)
		for r := posR; r < chamberWidth && r-posR < s.width; r++ {
			sr := r - posR
			if u < 0 || u >= len(c.rows) ||
				r < 0 || r >= chamberWidth ||
				su < shapeMax-s.height || su >= shapeMax ||
				sr < 0 || sr >= s.width {
				panic(fmt.Sprintf("range: u=%d, r=%d, su=%d, sr=%d, h=%d", u, r, su, sr, len(c.rows)))
			}
			if c.rows[u][r] && s.grid[su][sr] {
				return false // collision
			}
		}
	}
	return true
}

func (c *chamber) shapePlace(s shape, posU int, posR int, move int, drop uint64) {
	// add shape to chamber in position
	// add rows as needed
	add := posU + s.height - len(c.rows)
	if add > 0 {
		c.addRows(add)
	}
	// mark occupied shape area
	// for u := posU; u < len(c.rows) && u-posU < s.height; u++ {
	for u := posU; u-posU < s.height; u++ {
		su := shapeMax - (u - posU) - 1 // shape row is inverted (top/bottom)
		// for r := posR; r < chamberWidth && r-posR < s.width; r++ {
		for r := posR; r-posR < s.width; r++ {
			sr := r - posR
			if s.grid[su][sr] {
				c.rows[u][r] = true
			}
		}
		c.lastMove[u] = move
		c.lastDrop[u] = drop
	}
}

func (c *chamber) top() int {
	// find top of stack, rows are only added when used
	return len(c.rows)
}

// returns the number of rows repeated by the top line, or -1 if no pattern found
func (c *chamber) patternSearch() (int, uint64) {
	top := len(c.rows) - 1
	shapeCount := uint64(len(shapes))
	// min u value above 0 to ensure bottom doesn't impact patterns
	for u := len(c.rows) - 2; u >= top/2+10; u-- {
		mismatch := false
		// loop over size of pattern to ensure entire pattern is identical
		for offset := 0; offset < top-u && !mismatch; offset++ {
			// avoid identical drop for multi-row pieces, move and shape should be identical
			if c.lastDrop[u-offset] >= c.lastDrop[top-offset] || c.lastMove[u-offset] != c.lastMove[top-offset] || (c.lastDrop[u-offset]%shapeCount) != (c.lastDrop[top-offset]%shapeCount) {
				mismatch = true
			}
			for r := 0; r < chamberWidth && !mismatch; r++ {
				if c.rows[u-offset][r] != c.rows[top-offset][r] {
					mismatch = true
				}
			}
		}
		if !mismatch {
			return top - u, c.lastDrop[top] - c.lastDrop[u]
		}
	}
	return -1, 0
}

func (c *chamber) print() {
	// print from top to bottom
	for u := len(c.rows) - 1; u >= 0; u-- {
		fmt.Printf("  |")
		for r := 0; r < chamberWidth; r++ {
			if c.rows[u][r] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("|\n")
	}
	fmt.Printf("  +-------+\n")
}
