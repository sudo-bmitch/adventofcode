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
			'7': {X: 0, Y: 0},
			'8': {X: 0, Y: 1},
			'9': {X: 0, Y: 2},
			'4': {X: 1, Y: 0},
			'5': {X: 1, Y: 1},
			'6': {X: 1, Y: 2},
			'1': {X: 2, Y: 0},
			'2': {X: 2, Y: 1},
			'3': {X: 2, Y: 2},
			'0': {X: 3, Y: 1},
			'A': {X: 3, Y: 2},
		},
	}
	numPad.setValid()
	dirKeys := map[rune]grid.Pos{
		'^': {X: 0, Y: 1},
		'A': {X: 0, Y: 2},
		'<': {X: 1, Y: 0},
		'v': {X: 1, Y: 1},
		'>': {X: 1, Y: 2},
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
		// this isn't a random search, every valid path is shrinking either dx or dy
		dx, dy := end.X-start.X, end.Y-start.Y
		if dx != 0 {
			next := start
			if dx > 0 {
				next = next.MoveP(grid.Pos{X: 1})
			} else {
				next = next.MoveP(grid.Pos{X: -1})
			}
			if next == end {
				paths = append(paths, append(day21ClonePath(cur), next))
			} else if controller.valid[next] {
				queue = append(queue, append(day21ClonePath(cur), next))
			}
		}
		if dy != 0 {
			next := start
			if dy > 0 {
				next = next.MoveP(grid.Pos{Y: 1})
			} else {
				next = next.MoveP(grid.Pos{Y: -1})
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
			if path[i-1].X < path[i].X {
				entry += "v"
			} else if path[i-1].X > path[i].X {
				entry += "^"
			} else if path[i-1].Y < path[i].Y {
				entry += ">"
			} else if path[i-1].Y > path[i].Y {
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
