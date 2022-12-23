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

type elves struct {
	p                      map[pos]pos // cur to planned position
	xMin, xMax, yMin, yMax int
}

var search = [4][3]pos{
	{ // search north
		{x: -1, y: 0},
		{x: -1, y: -1},
		{x: -1, y: 1},
	},
	{ // search south
		{x: 1, y: 0},
		{x: 1, y: -1},
		{x: 1, y: 1},
	},
	{ // search west
		{x: 0, y: -1},
		{x: -1, y: -1},
		{x: 1, y: -1},
	},
	{ // search east
		{x: 0, y: 1},
		{x: -1, y: 1},
		{x: 1, y: 1},
	},
}

func main() {
	elvesCur := elves{p: map[pos]pos{}}
	x := 0
	in := bufio.NewScanner(os.Stdin)
	lineRE := regexp.MustCompile(`^[\.#]+$`)
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
		for y, c := range line {
			if c == '#' {
				elvesCur.set(pos{x: x, y: y})
			}
		}
		x++
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	for round := 0; round < 10; round++ {
		plannedPos := map[pos]int{}
		elvesNext := elves{p: map[pos]pos{}}
		// search until a new position found
		for cur := range elvesCur.p {
			plan := pos{x: cur.x, y: cur.y}
			found := true
			// check all adjacent squares
			for _, dx := range []int{-1, 0, 1} {
				for _, dy := range []int{-1, 0, 1} {
					if dx == 0 && dy == 0 { // not myself
						continue
					}
					if _, ok := elvesCur.p[pos{x: cur.x + dx, y: cur.y + dy}]; ok {
						found = false // another elf nearby, need to search
					}
				}
			}
			// look in the 4 directions
			for look := 0; look < 4 && !found; look++ {
				blocked := false
				// each search spans 3 points
				for span := 0; span < 3 && !blocked; span++ {
					searchPos := pos{
						x: cur.x + search[(look+round)%4][span].x,
						y: cur.y + search[(look+round)%4][span].y,
					}
					// look was blocked by a current elf
					if _, ok := elvesCur.p[searchPos]; ok {
						blocked = true
					}
				}
				// nothing blocking this search, found next location
				if !blocked {
					found = true
					plan = pos{
						x: cur.x + search[(look+round)%4][0].x,
						y: cur.y + search[(look+round)%4][0].y,
					}
				}
			}
			elvesCur.p[cur] = plan
			plannedPos[plan]++
		}
		// move if no other elves picked that position
		for cur, plan := range elvesCur.p {
			if plannedPos[plan] == 1 {
				elvesNext.set(plan)
			} else {
				elvesNext.set(cur)
			}
		}
		// switch to next position
		elvesCur = elvesNext
		if debug {
			fmt.Printf("debug: elves after round %d:\n", round+1)
			elvesCur.print()
		}
	}

	// count empty
	result := (elvesCur.xMax-elvesCur.xMin+1)*(elvesCur.yMax-elvesCur.yMin+1) - len(elvesCur.p)

	fmt.Printf("Result: %d\n", result)
}

func (e *elves) set(p pos) {
	e.p[pos{x: p.x, y: p.y}] = pos{x: p.x, y: p.y}
	if e.xMin == 0 && e.xMax == 0 && e.yMin == 0 && e.yMax == 0 {
		e.xMin, e.xMax, e.yMin, e.yMax = p.x, p.x, p.y, p.y
	} else {
		e.xMin = min(e.xMin, p.x)
		e.xMax = max(e.xMax, p.x)
		e.yMin = min(e.yMin, p.y)
		e.yMax = max(e.yMax, p.y)
	}
}

func (e *elves) print() {
	for x := e.xMin; x <= e.xMax; x++ {
		for y := e.yMin; y <= e.yMax; y++ {
			if _, ok := e.p[pos{x: x, y: y}]; ok {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// func mustAtoi(s string) int {
// 	i, err := strconv.Atoi(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return i
// }
