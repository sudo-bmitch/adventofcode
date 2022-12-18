package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type element int
type grid struct {
	g [][][]element
}

const (
	elUnknown element = iota
	elRock
	elSteam
)

func main() {
	g := grid{}
	in := bufio.NewScanner(os.Stdin)
	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		lineSplit := strings.Split(line, ",")
		if len(lineSplit) != 3 {
			fmt.Fprintf(os.Stderr, "invalid line: %s\n", line)
			return
		}
		x := mustAtoi(lineSplit[0])
		y := mustAtoi(lineSplit[1])
		z := mustAtoi(lineSplit[2])
		g.Set(x, y, z, elRock)
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	// scan for steam
	foundSteam := true
	for foundSteam {
		foundSteam = false
		for x := 0; x < len(g.g); x++ {
			for y := 0; y < len(g.g[x]); y++ {
				for z := 0; z < len(g.g[x][y]); z++ {
					// only parse unknown
					if g.g[x][y][z] != elUnknown {
						continue
					}
					// if any edge is and edge or steam, make this steam
					// 6 sides: x-1, x+1, y-1, y+1, z-1, z+1
					if x <= 0 || g.g[x-1][y][z] == elSteam ||
						x >= len(g.g)-1 || g.g[x+1][y][z] == elSteam ||
						y <= 0 || g.g[x][y-1][z] == elSteam ||
						y >= len(g.g[x])-1 || g.g[x][y+1][z] == elSteam ||
						z <= 0 || g.g[x][y][z-1] == elSteam ||
						z >= len(g.g[x][y])-1 || g.g[x][y][z+1] == elSteam {
						g.g[x][y][z] = elSteam
						foundSteam = true
					}
				}
			}
		}
	}

	// count exposed surface area
	count := 0
	for x := 0; x < len(g.g); x++ {
		for y := 0; y < len(g.g[x]); y++ {
			for z := 0; z < len(g.g[x][y]); z++ {
				// skip empty locations
				if g.g[x][y][z] != elRock {
					continue
				}
				// if set, and on edge or adjacent size unset
				// 6 sides: x-1, x+1, y-1, y+1, z-1, z+1
				if x <= 0 || g.g[x-1][y][z] == elSteam {
					count++
				}
				if x >= len(g.g)-1 || g.g[x+1][y][z] == elSteam {
					count++
				}
				if y <= 0 || g.g[x][y-1][z] == elSteam {
					count++
				}
				if y >= len(g.g[x])-1 || g.g[x][y+1][z] == elSteam {
					count++
				}
				if z <= 0 || g.g[x][y][z-1] == elSteam {
					count++
				}
				if z >= len(g.g[x][y])-1 || g.g[x][y][z+1] == elSteam {
					count++
				}
			}
		}
	}

	fmt.Printf("Result: %d\n", count)
}

func (g *grid) Set(x, y, z int, el element) {
	g.Expand(x, y, z)
	g.g[x][y][z] = el
}

func (g *grid) Expand(ex, ey, ez int) {
	// check if expand is needed
	if len(g.g) > ex && len(g.g[0]) > ey && len(g.g[0][0]) > ez {
		return
	}
	mx := max(len(g.g), ex+1)
	if mx > len(g.g) {
		g.g = append(g.g, make([][][]element, mx-len(g.g))...)
	}
	for x := 0; x < mx; x++ {
		my := max(len(g.g[0]), ey+1)
		if my > len(g.g[x]) {
			g.g[x] = append(g.g[x], make([][]element, my-len(g.g[x]))...)
		}
		for y := 0; y < my; y++ {
			mz := max(len(g.g[0][0]), ez+1)
			if mz > len(g.g[x][y]) {
				g.g[x][y] = append(g.g[x][y], make([]element, mz-len(g.g[x][y]))...)
			}
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
