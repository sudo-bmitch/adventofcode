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
	debug        = true
	dropMax      = 2022
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
	rows        [][chamberWidth]bool
	clearedRows int
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

	shapeCount := len(shapes)
	moveCount := len(moves)
	moveI := 0
	for drop := 0; drop < dropMax; drop++ {
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
				if debug && drop < 1 {
					fmt.Printf("Move %d\n", diffR)
				}
			}
			if !c.shapeAllowed(s, posU-1, posR) {
				break
			}
			posU--
			if debug && drop < 1 {
				fmt.Printf("Move down\n")
			}
		}
		c.shapePlace(s, posU, posR)
		if debug && (drop < 10 || (drop <= 1000 && drop%500 == 0)) {
			fmt.Printf("Dropped shape %d (height %d):\n", drop, c.top())
			c.print()
		}
		// if drop%100 == 0 {
		// 	c.clearUnreachable()
		// }
		if drop%100 == 0 {
			fmt.Printf("Drop %d, top=%d, cleared=%d\n", drop, c.top(), c.clearedRows)
		}
	}

	fmt.Printf("Result: %d\n", c.top()+c.clearedRows)
}

func (c *chamber) addRows(count int) {
	for i := 0; i < count; i++ {
		c.rows = append(c.rows, [chamberWidth]bool{})
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

func (c *chamber) shapePlace(s shape, posU int, posR int) {
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
	}
}

func (c *chamber) top() int {
	// find top of stack, rows are only added when used
	return len(c.rows)
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

func (c *chamber) clearUnreachable() {
	u := len(c.rows) - 1
	cursor := 0
	for u > 0 {
		if cursor > 0 && !c.rows[u][cursor] {
			cursor = 0
		}
		for cursor < chamberWidth && c.rows[u][cursor] {
			cursor++
		}
		if cursor == chamberWidth {
			break
		}
		u--
	}
	c.clearedRows += u + 1
	c.rows = c.rows[u+1:]
}
