package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	multiX = 4000000
	limitX = 4000000
	limitY = 4000000

	// limitX = 20
	// limitY = 20
)

type point struct{ x, y int }
type sensor struct {
	p point
	d int // distance
}

func main() {
	sensors := []sensor{}

	in := bufio.NewScanner(os.Stdin)
	lineRE := regexp.MustCompile(`^Sensor at x=([0-9\-]+), y=([0-9\-]+): closest beacon is at x=([0-9\-]+), y=([0-9\-]+)$`)
	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		match := lineRE.FindStringSubmatch(line)
		if match == nil || len(match) < 5 {
			fmt.Fprintf(os.Stderr, "failed parse line: %v\n", line)
			return
		}
		sP := point{
			x: mustAtoi(match[1]),
			y: mustAtoi(match[2]),
		}
		bP := point{
			x: mustAtoi(match[3]),
			y: mustAtoi(match[4]),
		}
		sensors = append(sensors, sensor{p: sP, d: mDistance(sP, bP)})
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}
	foundX, foundY := 0, 0
	fmt.Printf("Scanning...")
	for x := 0; x <= limitX; x++ {
		if x%1000 == 0 {
			fmt.Printf(".")
		}
		for y := 0; y <= limitY; y++ {
			inRange := false
			for _, sensor := range sensors {
				if mDistance(point{x: x, y: y}, sensor.p) <= sensor.d {
					// optimize by incrementing y to last point in range on row x
					y = sensor.p.y + sensor.d - diff(sensor.p.x, x)
					inRange = true
					break
				}
			}
			if !inRange {
				foundX, foundY = x, y
				fmt.Printf("\nFound gap at %d,%d: %d\n", x, y, x*multiX+y)
			}
		}
	}
	fmt.Printf("\n")
	fmt.Printf("Result: %d,%d = %d\n", foundX, foundY, foundX*multiX+foundY)
}

func mDistance(a, b point) int {
	return diff(a.x, b.x) + diff(a.y, b.y)
}

func diff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
