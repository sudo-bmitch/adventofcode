package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/sudo-bmitch/adventofcode/pkg/grid"
)

type day18Step struct {
	dirA, dirB grid.Dir
	lenA, lenB int
}

type day18Plot bool

func day18a(args []string, rdr io.Reader) (string, error) {
	minRow, minCol, maxRow, maxCol := 0, 0, 0, 0
	curPos := grid.Pos{Row: 0, Col: 0}
	plot := map[grid.Pos]day18Plot{
		curPos: true,
	}
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		s, err := day18Parse(line)
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for i := 0; i < s.lenA; i++ {
			curPos = curPos.MoveD(s.dirA)
			plot[curPos] = true
		}
		if curPos.Row < minRow {
			minRow = curPos.Row
		} else if curPos.Row > maxRow {
			maxRow = curPos.Row
		}
		if curPos.Col < minCol {
			minCol = curPos.Col
		} else if curPos.Col > maxCol {
			maxCol = curPos.Col
		}
	}
	sum := 0
	for r := minRow; r <= maxRow; r++ {
		corner := 0 // 1 corner entered from above, -1 corner entered from below
		verticals := 0
		for c := minCol; c <= maxCol; c++ {
			p := grid.Pos{Row: r, Col: c}
			if plot[p] {
				// look above and below to track corners and verticals
				above := plot[p.MoveD(grid.North)]
				below := plot[p.MoveD(grid.South)]
				if above && below {
					verticals++
				} else if above {
					if corner == 0 {
						corner = 1
					} else if corner == 1 {
						corner = 0
					} else if corner == -1 {
						corner = 0
						verticals++
					}
				} else if below {
					if corner == 0 {
						corner = -1
					} else if corner == -1 {
						corner = 0
					} else if corner == 1 {
						corner = 0
						verticals++
					}
				}
				sum++ // add currently dug space
			} else if corner != 0 {
				return "", fmt.Errorf("left a row to open ground at %s", p)
			} else if verticals%2 != 0 {
				sum++ // add space between even number of verticals (even/odd rule)
			}
		}
		if corner != 0 {
			return "", fmt.Errorf("corner did not return to 0 on row %d", r)
		}
		if verticals%2 != 0 {
			return "", fmt.Errorf("odd number of verticals on row %d: %d", r, verticals)
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day18b(args []string, rdr io.Reader) (string, error) {
	minRow, minCol, maxRow, maxCol := 0, 0, 0, 0
	curPos := grid.Pos{Row: 0, Col: 0}
	plot := map[grid.Pos]day18Plot{
		curPos: true,
	}
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		s, err := day18Parse(line)
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for i := 0; i < s.lenB; i++ {
			curPos = curPos.MoveD(s.dirB)
			plot[curPos] = true
		}
		if curPos.Row < minRow {
			minRow = curPos.Row
		} else if curPos.Row > maxRow {
			maxRow = curPos.Row
		}
		if curPos.Col < minCol {
			minCol = curPos.Col
		} else if curPos.Col > maxCol {
			maxCol = curPos.Col
		}
	}
	fmt.Fprintf(os.Stderr, "plot is from [%d,%d] to [%d,%d]\n", minRow, minCol, maxRow, maxCol)
	keys := make([][]int, maxRow-minRow+1)
	for p := range plot {
		kx := p.Row - minRow
		ky := p.Col - minCol
		keys[kx] = append(keys[kx], ky)
	}
	sum := 0
	for r := range keys {
		corner := 0 // 1 corner entered from above, -1 corner entered from below
		verticals := 0
		slices.Sort(keys[r])
		prevY := -1
		for _, c := range keys[r] {
			p := grid.Pos{Row: r + minRow, Col: c + minCol}
			// look above and below to track corners and verticals
			above := plot[p.MoveD(grid.North)]
			below := plot[p.MoveD(grid.South)]
			if above && below {
				verticals++
			} else if above {
				if corner == 0 {
					corner = 1
				} else if corner == 1 {
					corner = 0
				} else if corner == -1 {
					corner = 0
					verticals++
				}
			} else if below {
				if corner == 0 {
					corner = -1
				} else if corner == -1 {
					corner = 0
				} else if corner == 1 {
					corner = 0
					verticals++
				}
			}
			sum++ // add currently dug space
			if corner != 0 && verticals%2 != 0 {
				sum += c - prevY - 1 // add space between even number of verticals (even/odd rule)
			} else if corner == 0 && verticals%2 == 0 {
				sum += c - prevY - 1 // add space between even number of verticals (even/odd rule)
			}
			prevY = c
		}
		if corner != 0 {
			return "", fmt.Errorf("corner did not return to 0 on row %d", r+minRow)
		}
		if verticals%2 != 0 {
			return "", fmt.Errorf("odd number of verticals on row %d: %d", r+minRow, verticals)
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

var day18ParseRe = regexp.MustCompile(`^(U|D|L|R) ([0-9]+) \(#([0-9a-f]{6})\)$`)
var day18DirLookupA = map[string]grid.Dir{
	"U": grid.North,
	"D": grid.South,
	"L": grid.West,
	"R": grid.East,
}
var day18DirLookupB = map[rune]grid.Dir{
	'3': grid.North,
	'1': grid.South,
	'2': grid.West,
	'0': grid.East,
}

func day18Parse(line string) (day18Step, error) {
	match := day18ParseRe.FindStringSubmatch(line)
	if len(match) < 4 {
		return day18Step{}, fmt.Errorf("failed to parse line %s", line)
	}
	iA, err := strconv.Atoi(match[2])
	if err != nil {
		return day18Step{}, fmt.Errorf("failed to parse length %s on line %s: %w", match[2], line, err)
	}
	dA, ok := day18DirLookupA[match[1]]
	if !ok {
		return day18Step{}, fmt.Errorf("failed to lookup direction %s on line %s", match[1], line)
	}
	iB, err := strconv.ParseInt(match[3][:5], 16, 0)
	if err != nil {
		return day18Step{}, fmt.Errorf("failed to parse length B %s on line %s: %w", match[3], line, err)
	}
	dB, ok := day18DirLookupB[rune(match[3][5])]
	if !ok {
		return day18Step{}, fmt.Errorf("failed to lookup direction %s on line %s", match[3], line)
	}
	return day18Step{
		dirA: dA,
		lenA: iA,
		dirB: dB,
		lenB: int(iB),
	}, nil
}
