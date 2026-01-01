package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

func init() {
	registerDay("12a", day12a)
	registerDay("12b", day12b)
}

type day12Region struct {
	col, row int
	counts   []int
}

type day12CmpIn struct {
	piece1, piece2 int  // index of each piece
	dx, dy         int  // offset of piece 2 from piece 1
	r1, r2         int  // rotation of each piece (0-4)
	f1, f2         bool // flip of each piece (true/false)
}

func day12a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	debug := true
	pieces, regions, err := day12Parse(rdr)
	if err != nil {
		return "", err
	}
	// loop over the regions
	for regionNum, region := range regions {
		// TODO: potential to use a goroutine here
		// make a slice of placed pieces
		placed := []struct {
			piece   int
			x, y, r int
			f       bool
		}{}
		totalCount := 0
		for i, count := range region.counts {
			totalCount += count
			for range count {
				placed = append(placed, struct {
					piece int
					x     int
					y     int
					r     int
					f     bool
				}{
					piece: i,
				})
			}
		}
		if totalCount*9 < region.col*region.row {
			// shortcut
			sum++
			continue
		}
		curPlaced := 0
		fit := false
		for {
			fit = true
			if curPlaced > 0 {
				// check for any overlaps with previous pieces
				for prev := range curPlaced {
					if !day12Cmp(pieces, day12CmpIn{
						piece1: placed[prev].piece, piece2: placed[curPlaced].piece,
						dx: placed[curPlaced].x - placed[prev].x,
						dy: placed[curPlaced].y - placed[prev].y,
						r1: placed[prev].r,
						r2: placed[curPlaced].r,
						f1: placed[prev].f,
						f2: placed[curPlaced].f,
					}) {
						fit = false
						break
					}
				}
			}
			// if it fits, place the next piece
			if fit {
				// if debug {
				// 	fmt.Fprintf(os.Stderr, "region %d, placed %d pieces:", regionNum, curPlaced+1)
				// 	for i := range curPlaced + 1 {
				// 		fmt.Fprintf(os.Stderr, " [%d,%d,%d,%d,%t]", placed[i].piece, placed[i].x, placed[i].y, placed[i].r, placed[i].f)
				// 	}
				// 	fmt.Fprintf(os.Stderr, "\n")
				// }
				curPlaced++
				if curPlaced >= len(placed) {
					// last piece placed
					break
				}
				placed[curPlaced].r = 0
				placed[curPlaced].f = false
				if placed[curPlaced].piece == placed[curPlaced-1].piece {
					// same type of piece, start it after the prev
					placed[curPlaced].x = placed[curPlaced-1].x + 1
					placed[curPlaced].y = placed[curPlaced-1].y
					if placed[curPlaced].x+3 > region.col {
						// ran out of columns, move to next row
						placed[curPlaced].x = 0
						placed[curPlaced].y++
					}
					if placed[curPlaced].y+3 > region.row {
						// ran out of rows, this piece cannot fit, iterate previous piece
						fit = false
						curPlaced--
					}
				} else {
					// start the different type of pieces at the beginning
					placed[curPlaced].x = 0
					placed[curPlaced].y = 0
				}
			}
			if !fit {
				// iterate by rotating, flipping, and then moving each piece before falling back to previous piece and iterating that
				found := false
				for !found && curPlaced >= 0 {
					if !found {
						// try to rotate
						placed[curPlaced].r++
						if placed[curPlaced].r < 4 {
							found = true
						} else {
							placed[curPlaced].r = 0
						}
					}
					if !found {
						// try to flip
						placed[curPlaced].f = !placed[curPlaced].f
						if placed[curPlaced].f {
							found = true
						}
					}
					if !found {
						// try to shift
						placed[curPlaced].x++
						if placed[curPlaced].x+3 > region.col {
							placed[curPlaced].x = 0
							placed[curPlaced].y++
						}
						if placed[curPlaced].y+3 <= region.row {
							found = true
						}
					}
					if !found {
						// failed to iterate that piece, go back to previous piece
						curPlaced--
					}
				}
				if !found {
					// failed to find a solution
					break
				}
			}
		}
		if debug {
			fmt.Fprintf(os.Stderr, "finished checking region %d, fit is %t\n", regionNum, fit)
			if fit {
				fmt.Fprintf(os.Stderr, "result:")
				for _, p := range placed {
					fmt.Fprintf(os.Stderr, " [%d,%d,%d,%d,%t]", p.piece, p.x, p.y, p.r, p.f)
				}
				fmt.Fprintf(os.Stderr, "\n")
			}
		}
		if fit {
			sum++
		}
	}

	_, _ = pieces, regions
	return fmt.Sprintf("%d", sum), nil
}

func day12b(args []string, rdr io.Reader) (string, error) {
	sum := 0

	return fmt.Sprintf("%d", sum), nil
}

func day12Parse(rdr io.Reader) ([]grid.Grid, []day12Region, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return nil, nil, err
	}
	pieces := []grid.Grid{}
	regions := []day12Region{}
	inParts := strings.Split(strings.TrimSpace(string(in)), "\n\n")
	for i, pieceStr := range inParts[:len(inParts)-1] {
		pieceStr, ok := strings.CutPrefix(pieceStr, fmt.Sprintf("%d:\n", i))
		if !ok {
			return nil, nil, fmt.Errorf("piece missing expected prefix: %q", pieceStr)
		}
		piece, err := grid.FromString(pieceStr)
		if err != nil {
			return nil, nil, err
		}
		pieces = append(pieces, piece)
	}
	for line := range strings.SplitSeq(inParts[len(inParts)-1], "\n") {
		lineSplit := strings.SplitN(line, ": ", 2)
		if len(lineSplit) < 2 {
			return nil, nil, fmt.Errorf("did not find a colon space in %s", line)
		}
		colRow := strings.SplitN(lineSplit[0], "x", 2)
		if len(colRow) < 2 {
			return nil, nil, fmt.Errorf("did not find a col/row in line %s", line)
		}
		col, err := strconv.Atoi(colRow[0])
		if err != nil {
			return nil, nil, err
		}
		row, err := strconv.Atoi(colRow[1])
		if err != nil {
			return nil, nil, err
		}
		region := day12Region{
			col: col,
			row: row,
		}
		for numStr := range strings.SplitSeq(lineSplit[1], " ") {
			n, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, nil, err
			}
			region.counts = append(region.counts, n)
		}
		regions = append(regions, region)
	}
	return pieces, regions, nil
}

var day12CmpMemoize = map[day12CmpIn]bool{}

func day12Cmp(pieces []grid.Grid, cmpIn day12CmpIn) bool {
	// pieces are all 3x3, no overlaps when dx or dy exceeds that
	if cmpIn.dx > 3 || cmpIn.dy > 3 || cmpIn.dx < -3 || cmpIn.dy < -3 {
		return true
	}
	if result, ok := day12CmpMemoize[cmpIn]; ok {
		return result
	}
	// put each piece on a grid
	p1 := [3][3]bool{}
	for pos, r := range pieces[cmpIn.piece1].Walk() {
		if r == '.' {
			continue
		} else if r != '#' {
			panic("invalid piece")
		}
		if pos.Col < 0 || pos.Col >= 3 || pos.Row < 0 || pos.Row >= 3 {
			panic("invalid piece")
		}
		p1[pos.Col][pos.Row] = true
	}
	p2 := [3][3]bool{}
	for pos, r := range pieces[cmpIn.piece2].Walk() {
		if r == '.' {
			continue
		} else if r != '#' {
			panic("invalid piece")
		}
		if pos.Col < 0 || pos.Col >= 3 || pos.Row < 0 || pos.Row >= 3 {
			panic("invalid piece")
		}
		p2[pos.Col][pos.Row] = true
	}
	// flip
	if cmpIn.f1 {
		p1[0][0], p1[2][0] = p1[2][0], p1[0][0]
		p1[0][1], p1[2][1] = p1[2][0], p1[0][1]
		p1[0][2], p1[2][2] = p1[2][0], p1[0][2]
	}
	if cmpIn.f2 {
		p2[0][0], p2[2][0] = p2[2][0], p2[0][0]
		p2[0][1], p2[2][1] = p2[2][0], p2[0][1]
		p2[0][2], p2[2][2] = p2[2][0], p2[0][2]
	}
	// rotate
	for range cmpIn.r1 {
		p1[0][0], p1[2][0], p1[2][2], p1[0][2] = p1[2][0], p1[2][2], p1[0][2], p1[0][0]
		p1[1][0], p1[2][1], p1[1][2], p1[0][1] = p1[2][1], p1[1][2], p1[0][1], p1[1][0]
	}
	for range cmpIn.r2 {
		p2[0][0], p2[2][0], p2[2][2], p2[0][2] = p2[2][0], p2[2][2], p2[0][2], p2[0][0]
		p2[1][0], p2[2][1], p2[1][2], p2[0][1] = p2[2][1], p2[1][2], p2[0][1], p2[1][0]
	}
	// check for overlaps
	result := true
	for x2 := range 3 {
		for y2 := range 3 {
			x1 := x2 - cmpIn.dx
			y1 := y2 - cmpIn.dy
			if x1 < 0 || x1 > 2 || y1 < 0 || y1 > 2 {
				continue
			}
			if p1[x1][y1] && p2[x2][y2] {
				result = false
				break
			}
		}
		if !result {
			break
		}
	}
	// store result
	day12CmpMemoize[cmpIn] = result
	return result
}
