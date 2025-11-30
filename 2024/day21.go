package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("21a", day21a)
	registerDay("21b", day21b)
}

func day21a(args []string, rdr io.Reader) (string, error) {
	sum, err := day21Run(rdr, 2)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", sum), nil
}

func day21b(args []string, rdr io.Reader) (string, error) {
	sum, err := day21Run(rdr, 25)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", sum), nil
}

func day21Run(rdr io.Reader, numControllers int) (int, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(in), "\n")

	// define controllers and link them
	numPad := day21Controller{
		keys: map[rune]grid.Pos{
			'7': {Row: 0, Col: 0},
			'8': {Row: 0, Col: 1},
			'9': {Row: 0, Col: 2},
			'4': {Row: 1, Col: 0},
			'5': {Row: 1, Col: 1},
			'6': {Row: 1, Col: 2},
			'1': {Row: 2, Col: 0},
			'2': {Row: 2, Col: 1},
			'3': {Row: 2, Col: 2},
			'0': {Row: 3, Col: 1},
			'A': {Row: 3, Col: 2},
		},
	}
	numPad.setValid()
	dirKeys := map[rune]grid.Pos{
		'^': {Row: 0, Col: 1},
		'A': {Row: 0, Col: 2},
		'<': {Row: 1, Col: 0},
		'v': {Row: 1, Col: 1},
		'>': {Row: 1, Col: 2},
	}
	dirPads := make([]day21Controller, numControllers)
	for i := range dirPads {
		dirPads[i].keys = dirKeys
		dirPads[i].setValid()
		if i > 0 {
			dirPads[i-1].parent = &dirPads[i]
		}
	}
	numPad.parent = &dirPads[0]

	// iterate over inputs to find length of each code
	sum := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		best := numPad.findBest(line)
		fmt.Printf("Code %s solved with %d keys\n", line, best)
		code, err := strconv.Atoi(line[:len(line)-1])
		if err != nil {
			return 0, err
		}
		sum += code * best
	}
	return sum, nil
}

type day21Controller struct {
	keys   map[rune]grid.Pos // position of each key on this pad
	valid  map[grid.Pos]bool // reverse key map to flag invalid positions
	cache  map[string]int    // cache tracks the best string for this controller
	parent *day21Controller  // parent references the next level controller, nil for human
}

func (c *day21Controller) setValid() {
	c.valid = map[grid.Pos]bool{}
	for _, p := range c.keys {
		c.valid[p] = true
	}
}

// findBest runs on a set of keys, always beginning on A, which should also be the last key in every sequence
func (c *day21Controller) findBest(keys string) int {
	if c.cache == nil {
		c.cache = map[string]int{}
	}
	if result, ok := c.cache[keys]; ok {
		return result
	}
	// start at A key
	cur := c.keys['A']
	result := 0
	// go through each key tracking cur location
	for _, key := range keys {
		// generate a list of options from cur to desired key
		end := c.keys[key]
		opts := day21GenOpts(*c, cur, end)
		best := 0
		if c.parent != nil {
			// if parent controller exists, run through parent controller to find shortest length for the list of options
			parentOptLens := []int{}
			for _, opt := range opts {
				parentOptLens = append(parentOptLens, c.parent.findBest(opt))
			}
			best = parentOptLens[0]
			for _, b := range parentOptLens {
				if b < best {
					best = b
				}
			}
		} else {
			// pick shortest entry and append to result
			best = len(opts[0])
			for _, opt := range opts {
				if len(opt) < best {
					best = len(opt)
				}
			}
		}
		result += best
		// update cur location for next key
		cur = end
	}
	// cache result
	c.cache[keys] = result
	return result
}

// day21GenOpts creates a list of direction key options that will go from the start to end position
func day21GenOpts(controller day21Controller, start, end grid.Pos) []string {
	if start == end {
		return []string{"A"}
	}
	queue := [][]grid.Pos{{start}}
	// each path is a list of coordinates from start to end
	paths := [][]grid.Pos{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		start := cur[len(cur)-1]
		// this isn't a random search, every valid path is shrinking either dr or dc
		dr, dc := end.Row-start.Row, end.Col-start.Col
		if dr != 0 {
			next := start
			if dr > 0 {
				next = next.MoveP(grid.Pos{Row: 1})
			} else {
				next = next.MoveP(grid.Pos{Row: -1})
			}
			if next == end {
				paths = append(paths, append(day21ClonePath(cur), next))
			} else if controller.valid[next] {
				queue = append(queue, append(day21ClonePath(cur), next))
			}
		}
		if dc != 0 {
			next := start
			if dc > 0 {
				next = next.MoveP(grid.Pos{Col: 1})
			} else {
				next = next.MoveP(grid.Pos{Col: -1})
			}
			if next == end {
				paths = append(paths, append(cur, next))
			} else if controller.valid[next] {
				queue = append(queue, append(cur, next))
			}
		}
	}
	// convert paths into directions depending on dx or dy change per hop
	list := []string{}
	for _, path := range paths {
		entry := ""
		for i := 1; i < len(path); i++ {
			if path[i-1].Row < path[i].Row {
				entry += "v"
			} else if path[i-1].Row > path[i].Row {
				entry += "^"
			} else if path[i-1].Col < path[i].Col {
				entry += ">"
			} else if path[i-1].Col > path[i].Col {
				entry += "<"
			}
		}
		// every path ends with an A to press the button
		entry += "A"
		list = append(list, entry)
	}
	return list
}

func day21ClonePath(orig []grid.Pos) []grid.Pos {
	clone := make([]grid.Pos, len(orig))
	copy(clone, orig)
	return clone
}
