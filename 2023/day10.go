package main

import (
	"fmt"
	"io"
	"strings"
)

type day10Pos struct {
	x, y int
}

var (
	day10DNorth = day10Pos{x: -1}
	day10DEast  = day10Pos{y: 1}
	day10DSouth = day10Pos{x: 1}
	day10DWest  = day10Pos{y: -1}
)

func (p day10Pos) step(d day10Pos) day10Pos {
	return day10Pos{
		x: p.x + d.x,
		y: p.y + d.y,
	}
}

type day10Map struct {
	start day10Pos
	pipes [][]rune
	max   day10Pos
}

type day10Dir struct {
	north, east, south, west bool
}

type day10Side struct {
	pipe                           bool
	north, east, south, west, open int
}

const (
	day10SideR = 1
	day10SideL = 2

	day10North = 1
	day10East  = 2
	day10South = 3
	day10West  = 4
)

var day10Connect = map[rune]day10Dir{
	'|': {north: true, south: true},
	'-': {west: true, east: true},
	'L': {north: true, east: true},
	'J': {north: true, west: true},
	'7': {west: true, south: true},
	'F': {east: true, south: true},
	'.': {},
	'S': {north: true, east: true, south: true, west: true},
}

func day10a(args []string, rdr io.Reader) (string, error) {
	max := 0
	m, err := day10Parse(rdr)
	if err != nil {
		return "", err
	}
	visited := map[day10Pos]int{
		m.start: 0,
	}
	search := []day10Pos{m.start}
	for len(search) > 0 {
		nextSearch := make([]day10Pos, 0, len(search))
		max++
		for _, s := range search {
			trySearch := m.next(s)
			for _, t := range trySearch {
				if _, ok := visited[t]; !ok {
					visited[t] = max
					nextSearch = append(nextSearch, t)
				}
			}
		}
		search = nextSearch
	}
	max-- // remove last increment
	return fmt.Sprintf("%d", max), nil
}

func day10b(args []string, rdr io.Reader) (string, error) {
	m, err := day10Parse(rdr)
	if err != nil {
		return "", err
	}
	sides := make([][]day10Side, len(m.pipes))
	y := len(m.pipes[0])
	for x := range sides {
		sides[x] = make([]day10Side, y)
	}
	// determine the shape of the start, pick an exit, and set the side values
	cur := m.start
	dirs := day10Dir{}
	if next := cur.step(day10DNorth); m.posValid(next) && day10Connect[m.pipes[next.x][next.y]].south {
		dirs.north = true
	}
	if next := cur.step(day10DEast); m.posValid(next) && day10Connect[m.pipes[next.x][next.y]].west {
		dirs.east = true
	}
	if next := cur.step(day10DSouth); m.posValid(next) && day10Connect[m.pipes[next.x][next.y]].north {
		dirs.south = true
	}
	if next := cur.step(day10DWest); m.posValid(next) && day10Connect[m.pipes[next.x][next.y]].east {
		dirs.west = true
	}
	var enter int
	var next day10Pos
	switch {
	case dirs.north && dirs.east:
		sides[cur.x][cur.y] = day10Side{pipe: true, west: day10SideL, south: day10SideL}
		next = cur.step(day10DNorth)
		enter = day10South
	case dirs.north && dirs.south:
		sides[cur.x][cur.y] = day10Side{pipe: true, west: day10SideL, east: day10SideR}
		next = cur.step(day10DNorth)
		enter = day10South
	case dirs.north && dirs.west:
		sides[cur.x][cur.y] = day10Side{pipe: true, east: day10SideR, south: day10SideR}
		next = cur.step(day10DNorth)
		enter = day10South
	case dirs.east && dirs.south:
		sides[cur.x][cur.y] = day10Side{pipe: true, north: day10SideL, west: day10SideL}
		next = cur.step(day10DEast)
		enter = day10West
	case dirs.east && dirs.west:
		sides[cur.x][cur.y] = day10Side{pipe: true, north: day10SideL, south: day10SideR}
		next = cur.step(day10DEast)
		enter = day10West
	case dirs.south && dirs.west:
		sides[cur.x][cur.y] = day10Side{pipe: true, north: day10SideL, east: day10SideL}
		next = cur.step(day10DSouth)
		enter = day10North
	default:
		return "", fmt.Errorf("unhandled start scenario: %v", dirs)
	}

	// walk until returning to the start
	cur = next
	for cur != m.start {
		// verity entry dir
		if !m.posValid(cur) ||
			(enter == day10North && !day10Connect[m.pipes[cur.x][cur.y]].north) ||
			(enter == day10East && !day10Connect[m.pipes[cur.x][cur.y]].east) ||
			(enter == day10South && !day10Connect[m.pipes[cur.x][cur.y]].south) ||
			(enter == day10West && !day10Connect[m.pipes[cur.x][cur.y]].west) {
			return "", fmt.Errorf("invalid position (%v) or entrypoint (%d) for rune \"%c\"", cur, enter, m.pipes[cur.x][cur.y])
		}
		// identify exit dir, set sides for cur, update enter and cur for next position
		switch {
		case enter == day10North && day10Connect[m.pipes[cur.x][cur.y]].east:
			sides[cur.x][cur.y] = day10Side{pipe: true, south: day10SideR, west: day10SideR}
			enter = day10West
			cur = cur.step(day10DEast)
		case enter == day10North && day10Connect[m.pipes[cur.x][cur.y]].south:
			sides[cur.x][cur.y] = day10Side{pipe: true, east: day10SideL, west: day10SideR}
			enter = day10North
			cur = cur.step(day10DSouth)
		case enter == day10North && day10Connect[m.pipes[cur.x][cur.y]].west:
			sides[cur.x][cur.y] = day10Side{pipe: true, east: day10SideL, south: day10SideL}
			enter = day10East
			cur = cur.step(day10DWest)
		case enter == day10East && day10Connect[m.pipes[cur.x][cur.y]].south:
			sides[cur.x][cur.y] = day10Side{pipe: true, north: day10SideR, west: day10SideR}
			enter = day10North
			cur = cur.step(day10DSouth)
		case enter == day10East && day10Connect[m.pipes[cur.x][cur.y]].west:
			sides[cur.x][cur.y] = day10Side{pipe: true, north: day10SideR, south: day10SideL}
			enter = day10East
			cur = cur.step(day10DWest)
		case enter == day10East && day10Connect[m.pipes[cur.x][cur.y]].north:
			sides[cur.x][cur.y] = day10Side{pipe: true, south: day10SideL, west: day10SideL}
			enter = day10South
			cur = cur.step(day10DNorth)
		case enter == day10South && day10Connect[m.pipes[cur.x][cur.y]].west:
			sides[cur.x][cur.y] = day10Side{pipe: true, north: day10SideR, east: day10SideR}
			enter = day10East
			cur = cur.step(day10DWest)
		case enter == day10South && day10Connect[m.pipes[cur.x][cur.y]].north:
			sides[cur.x][cur.y] = day10Side{pipe: true, east: day10SideR, west: day10SideL}
			enter = day10South
			cur = cur.step(day10DNorth)
		case enter == day10South && day10Connect[m.pipes[cur.x][cur.y]].east:
			sides[cur.x][cur.y] = day10Side{pipe: true, west: day10SideL, north: day10SideL}
			enter = day10West
			cur = cur.step(day10DEast)
		case enter == day10West && day10Connect[m.pipes[cur.x][cur.y]].north:
			sides[cur.x][cur.y] = day10Side{pipe: true, east: day10SideR, south: day10SideR}
			enter = day10South
			cur = cur.step(day10DNorth)
		case enter == day10West && day10Connect[m.pipes[cur.x][cur.y]].east:
			sides[cur.x][cur.y] = day10Side{pipe: true, north: day10SideL, south: day10SideR}
			enter = day10West
			cur = cur.step(day10DEast)
		case enter == day10West && day10Connect[m.pipes[cur.x][cur.y]].south:
			sides[cur.x][cur.y] = day10Side{pipe: true, north: day10SideL, east: day10SideL}
			enter = day10North
			cur = cur.step(day10DSouth)
		default:
		}
	}

	// fill in open spaces by scanning for non-pipe spaces and claim each for a side until no more claims made
	unknown := true // track if there are unknown spaces remaining to discover
	for unknown {
		unknown = false
		for x := range sides {
			for y := range sides[x] {
				if !sides[x][y].pipe && sides[x][y].open == 0 {
					cur = day10Pos{x: x, y: y}
					if try := cur.step(day10DNorth); m.posValid(try) {
						if sides[try.x][try.y].pipe {
							sides[x][y].open = sides[try.x][try.y].south
						} else if sides[try.x][try.y].open != 0 {
							sides[x][y].open = sides[try.x][try.y].open
						}
					}
					if try := cur.step(day10DEast); m.posValid(try) && sides[x][y].open == 0 {
						if sides[try.x][try.y].pipe {
							sides[x][y].open = sides[try.x][try.y].west
						} else if sides[try.x][try.y].open != 0 {
							sides[x][y].open = sides[try.x][try.y].open
						}
					}
					if try := cur.step(day10DSouth); m.posValid(try) && sides[x][y].open == 0 {
						if sides[try.x][try.y].pipe {
							sides[x][y].open = sides[try.x][try.y].north
						} else if sides[try.x][try.y].open != 0 {
							sides[x][y].open = sides[try.x][try.y].open
						}
					}
					if try := cur.step(day10DWest); m.posValid(try) && sides[x][y].open == 0 {
						if sides[try.x][try.y].pipe {
							sides[x][y].open = sides[try.x][try.y].east
						} else if sides[try.x][try.y].open != 0 {
							sides[x][y].open = sides[try.x][try.y].open
						}
					}
					if sides[x][y].open == 0 {
						unknown = true
					}
				}
			}
		}
	}
	// identify invalid side by checking 0,0 and determine if it's L or R
	var outside int
	if sides[0][0].pipe {
		outside = sides[0][0].north
	} else if sides[0][0].open != 0 {
		outside = sides[0][0].open
	} else {
		return "", fmt.Errorf("unable to determine outside from %v", sides[0][0])
	}

	// scan again to count all open spaces for the valid side
	sum := 0
	for x := range sides {
		for y := range sides[x] {
			if sides[x][y].open != 0 && sides[x][y].open != outside {
				sum++
			}
		}
	}

	return fmt.Sprintf("%d", sum), nil
}

func day10Parse(rdr io.Reader) (day10Map, error) {
	m := day10Map{
		pipes: [][]rune{},
	}
	lines, err := io.ReadAll(rdr)
	if err != nil {
		return m, err
	}
	for x, line := range strings.Split(string(lines), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if m.max.y == 0 {
			m.max.y = len(line)
		} else if m.max.y != len(line) {
			return m, fmt.Errorf("line length mismatch %d != %d", m.max.y, len(line))
		}
		m.pipes = append(m.pipes, []rune(line))
		for y, r := range []rune(line) {
			if r == 'S' {
				m.start = day10Pos{x: x, y: y}
			}
		}
	}
	m.max.x = len(m.pipes)
	return m, nil
}

func (m day10Map) posValid(cur day10Pos) bool {
	if cur.x < 0 || cur.y < 0 || cur.x >= m.max.x || cur.y >= m.max.y {
		return false
	}
	return true
}

func (m *day10Map) next(cur day10Pos) []day10Pos {
	next := make([]day10Pos, 0, 2)
	if cur.x < 0 || cur.y < 0 || cur.x >= len(m.pipes) || cur.y >= len(m.pipes[cur.x]) {
		return next
	}
	conn := day10Connect[m.pipes[cur.x][cur.y]]
	if conn.north && cur.x > 0 {
		nextTry := day10Pos{x: cur.x - 1, y: cur.y}
		if day10Connect[m.pipes[nextTry.x][nextTry.y]].south {
			next = append(next, nextTry)
		}
	}
	if conn.east && cur.y < len(m.pipes[cur.x])-1 {
		nextTry := day10Pos{x: cur.x, y: cur.y + 1}
		if day10Connect[m.pipes[nextTry.x][nextTry.y]].west {
			next = append(next, nextTry)
		}
	}
	if conn.south && cur.x < len(m.pipes)-1 {
		nextTry := day10Pos{x: cur.x + 1, y: cur.y}
		if day10Connect[m.pipes[nextTry.x][nextTry.y]].north {
			next = append(next, nextTry)
		}
	}
	if conn.west && cur.y > 0 {
		nextTry := day10Pos{x: cur.x, y: cur.y - 1}
		if day10Connect[m.pipes[nextTry.x][nextTry.y]].east {
			next = append(next, nextTry)
		}
	}
	return next
}
