package main

import (
	"fmt"
	"io"
	"maps"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("09a", day09a)
	registerDay("09b", day09b)
}

func day09a(args []string, rdr io.Reader) (string, error) {
	largest := 0
	debug := false
	posList := []grid.Pos{}

	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	for line := range strings.SplitSeq(strings.TrimSpace(string(in)), "\n") {
		xy := strings.SplitN(line, ",", 2)
		if len(xy) != 2 {
			return "", fmt.Errorf("failed to split line %q", line)
		}
		x, err := strconv.Atoi(xy[0])
		if err != nil {
			return "", fmt.Errorf("failed to convert number %s from line %q: %w", xy[0], line, err)
		}
		y, err := strconv.Atoi(xy[1])
		if err != nil {
			return "", fmt.Errorf("failed to convert number %s from line %q: %w", xy[1], line, err)
		}
		posList = append(posList, grid.Pos{Col: x, Row: y})
	}
	for i := range posList {
		for j := i + 1; j < len(posList); j++ {
			area := (day09Abs(posList[i].Col-posList[j].Col) + 1) * (day09Abs(posList[i].Row-posList[j].Row) + 1)
			if area < 0 {
				area *= -1
			}
			if area > largest {
				largest = area
				if debug {
					fmt.Fprintf(os.Stderr, "found larger area %d between %d,%d and %d,%d\n", area, posList[i].Col, posList[i].Row, posList[j].Col, posList[j].Row)
				}
			}
		}
	}
	return fmt.Sprintf("%d", largest), nil
}

func day09b(args []string, rdr io.Reader) (string, error) {
	largest := 0
	debug := false

	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	minX, maxX, minY, maxY := math.MaxInt, 0, math.MaxInt, 0
	posList := []grid.Pos{}
	for line := range strings.SplitSeq(strings.TrimSpace(string(in)), "\n") {
		xy := strings.SplitN(line, ",", 2)
		if len(xy) != 2 {
			return "", fmt.Errorf("failed to split line %q", line)
		}
		x, err := strconv.Atoi(xy[0])
		if err != nil {
			return "", fmt.Errorf("failed to convert number %s from line %q: %w", xy[0], line, err)
		}
		y, err := strconv.Atoi(xy[1])
		if err != nil {
			return "", fmt.Errorf("failed to convert number %s from line %q: %w", xy[1], line, err)
		}
		minX = min(minX, x)
		maxX = max(maxX, x)
		minY = min(minY, y)
		maxY = max(maxY, y)
		posList = append(posList, grid.Pos{Col: x, Row: y})
	}

	type day09RowCross struct {
		north, south bool
	}
	rowCrosses := map[int]map[int]*day09RowCross{}
	for i, p1 := range posList {
		i2 := (i + 1) % len(posList)
		p2 := posList[i2]
		north := false
		if p1.Row == p2.Row {
			continue // nothing to track on horizontal
		}
		if p1.Col != p2.Col {
			return "", fmt.Errorf("dir not straight between %v and %v", p1, p2)
		}
		if p1.Row > p2.Row {
			north = true
		}
		cur := p1
		for cur != p2 {
			if rowCrosses[cur.Row] == nil {
				rowCrosses[cur.Row] = map[int]*day09RowCross{}
			}
			if rowCrosses[cur.Row][cur.Col] == nil {
				rowCrosses[cur.Row][cur.Col] = &day09RowCross{}
			}
			// marking the side of the exit, and move to the next square
			if north {
				rowCrosses[cur.Row][cur.Col].north = true
				cur = cur.MoveD(grid.North)
			} else {
				rowCrosses[cur.Row][cur.Col].south = true
				cur = cur.MoveD(grid.South)
			}
			if rowCrosses[cur.Row] == nil {
				rowCrosses[cur.Row] = map[int]*day09RowCross{}
			}
			if rowCrosses[cur.Row][cur.Col] == nil {
				rowCrosses[cur.Row][cur.Col] = &day09RowCross{}
			}
			// mark the side we entered
			if north {
				rowCrosses[cur.Row][cur.Col].south = true
			} else {
				rowCrosses[cur.Row][cur.Col].north = true
			}
		}
	}
	debugGrid := false
	if debugGrid {
		fmt.Fprintf(os.Stderr, "Filled grid:\n")
		for row := range maxY + 2 {
			for col := range maxX + 2 {
				if rowCrosses[row] != nil && rowCrosses[row][col] != nil {
					if rowCrosses[row][col].north && rowCrosses[row][col].south {
						fmt.Fprintf(os.Stderr, "|")
					} else if rowCrosses[row][col].north {
						fmt.Fprintf(os.Stderr, "N")
					} else {
						fmt.Fprintf(os.Stderr, "S")
					}
				} else {
					fmt.Fprintf(os.Stderr, ".")
				}
			}
			fmt.Fprintf(os.Stderr, "\n")
		}
	}

	for i := range posList {
		if debug {
			fmt.Fprintf(os.Stderr, "[i = %d] checking rectangles from corner [%d,%d]\n", i, posList[i].Col, posList[i].Row)
		}
		for j := i + 1; j < len(posList); j++ {
			area := (day09Abs(posList[i].Col-posList[j].Col) + 1) * (day09Abs(posList[i].Row-posList[j].Row) + 1)
			if area > largest {
				// check if all points are within the lines
				inside := true
				colStart := min(posList[i].Col, posList[j].Col)
				colStop := max(posList[i].Col, posList[j].Col)
				for row := min(posList[i].Row, posList[j].Row); inside && row <= max(posList[i].Row, posList[j].Row); row++ {
					colLast := -1
					countN, countS := 0, 0
					for _, col := range slices.Sorted(maps.Keys(rowCrosses[row])) {
						// previously outside, this col will put me on a corner or line,
						// verify the gap between colLast and col is > 0 and then
						// check colLast+1 to col-1 to make sure none of those values are inside the rectangle
						if countN%2 == 0 && countS%2 == 0 && colLast+1 <= colStop && colStart <= col-1 && col-colLast-1 > 0 {
							inside = false
							break
						}
						// then move col tracking the crossing
						if rowCrosses[row][col].north {
							countN++
						}
						if rowCrosses[row][col].south {
							countS++
						}
						colLast = col
					}
					// last crossing should have put us back outside, ensure colLast+1 is not inside the rectangle
					if inside && (countN%2 != 0 || countS%2 != 0) {
						return "", fmt.Errorf("bad exit from row %d, countN=%d, countS=%d, rowCrosses[row]=%v", row, countN, countS, rowCrosses[row])
					}
					if colLast+1 <= colStop {
						inside = false
					}
				}
				if inside {
					largest = area
					if debug {
						fmt.Fprintf(os.Stderr, "[i = %d] found larger area %d between [%d,%d] and [%d,%d]\n", i, area, posList[i].Col, posList[i].Row, posList[j].Col, posList[j].Row)
					}
				}
			}
		}
	}
	return fmt.Sprintf("%d", largest), nil
	// 1516172795
}

func day09Abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}
