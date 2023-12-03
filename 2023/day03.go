package main

import (
	"bufio"
	"fmt"
	"io"
)

type grid struct {
	numbers map[pos]num
	objs    map[pos]rune
}

type pos struct {
	x, y int
}

type num struct {
	value int
	size  int
}

func day03a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	g, err := day03Parse(rdr)
	if err != nil {
		return "", err
	}
	for p, n := range g.numbers {
		if day03ObjNeighbor(g, n, p) {
			sum += n.value
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day03b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	g, err := day03Parse(rdr)
	if err != nil {
		return "", err
	}
	for p, o := range g.objs {
		if o == '*' {
			nList := day03NumNeighbors(g, p)
			if len(nList) > 2 {
				return "", fmt.Errorf("found more than 2 neighbors at x=%d, y=%d", p.x, p.y)
			}
			if len(nList) == 2 {
				ratio := nList[0].value * nList[1].value
				sum += ratio
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

func day03ObjNeighbor(g grid, n num, p pos) bool {
	for x := p.x - 1; x <= p.x+1; x++ {
		if x < 0 {
			continue
		}
		for y := p.y - n.size; y <= p.y+1; y++ {
			if y < 0 {
				continue
			}
			if _, ok := g.objs[pos{x: x, y: y}]; ok {
				return true
			}
		}
	}
	return false
}

func day03NumNeighbors(g grid, p pos) []num {
	nList := []num{}
	for pCur, nCur := range g.numbers {
		if p.x >= pCur.x-1 && p.x <= pCur.x+1 && p.y >= pCur.y-nCur.size && p.y <= pCur.y+1 {
			nList = append(nList, nCur)
		}
	}
	return nList
}

func day03Parse(rdr io.Reader) (grid, error) {
	g := grid{
		numbers: map[pos]num{},
		objs:    map[pos]rune{},
	}
	x := 0
	n := 0
	nSize := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		for y, c := range line {
			switch c {
			case '.':
				if nSize > 0 {
					g.numbers[pos{x: x, y: y - 1}] = num{value: n, size: nSize}
					n = 0
					nSize = 0
				}
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				n = n*10 + int(c-'0')
				nSize++
			default:
				g.objs[pos{x: x, y: y}] = c
				if nSize > 0 {
					g.numbers[pos{x: x, y: y - 1}] = num{value: n, size: nSize}
					n = 0
					nSize = 0
				}
			}
		}
		// handle number on the far right of grid
		if nSize > 0 {
			g.numbers[pos{x: x, y: len(line) - 1}] = num{value: n, size: nSize}
			n = 0
			nSize = 0
		}
		x++
	}
	return g, nil
}
