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

// const targetY = 10
const targetY = 2000000

type span struct {
	start, end int
}
type spans []span

func (sl *spans) addSpan(s span) {
	if s.start > s.end {
		panic(fmt.Sprintf("invalid span: start %d, end %d", s.start, s.end))
	}
	// insert where start is before the current entry start
	insertPos := 0
	for insertPos < len(*sl) && s.start > (*sl)[insertPos].end {
		insertPos++
	}
	if insertPos == len(*sl) {
		fmt.Printf("addSpan: added to end: %d - %d\n", s.start, s.end)
		(*sl) = append((*sl), s)
	} else {
		fmt.Printf("addSpan: inserted at %d: %d (%d) - %d\n", insertPos, s.start, min(s.start, (*sl)[insertPos].start), s.end)
		s.start = min(s.start, (*sl)[insertPos].start)
		(*sl) = append((*sl)[:insertPos+1], (*sl)[insertPos:]...)
		(*sl)[insertPos] = s
	}
	// if next start is less than current end, merge two entries, and repeat
	for insertPos+1 < len(*sl) && (*sl)[insertPos].end >= (*sl)[insertPos+1].start {
		fmt.Printf("addSpan: deleting after %d: %d - %d, new end %d\n", insertPos, (*sl)[insertPos+1].start, (*sl)[insertPos+1].end, max((*sl)[insertPos].end, (*sl)[insertPos+1].end))
		(*sl)[insertPos].end = max((*sl)[insertPos].end, (*sl)[insertPos+1].end)
		(*sl) = append((*sl)[:insertPos+1], (*sl)[insertPos+2:]...)
	}
}

func main() {
	objects := []int{}
	rowSpans := spans{}
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
		sX := mustAtoi(match[1])
		sY := mustAtoi(match[2])
		bX := mustAtoi(match[3])
		bY := mustAtoi(match[4])
		dBeacon := mDistance(sX, sY, bX, bY)
		dTarget := mDistance(sX, sY, sX, targetY)
		if dTarget <= dBeacon {
			width := dBeacon - dTarget
			fmt.Printf("Line overlap: %s, width %d, sX %d\n", line, width, sX)
			rowSpans.addSpan(span{start: sX - width, end: sX + width})
		} else {
			fmt.Printf("Line ignored: %s, dBeacon %d, dTarget %d\n", line, dBeacon, dTarget)
		}
		if sY == targetY {
			objects = append(objects, sX)
		}
		if bY == targetY {
			objects = append(objects, bX)
		}
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}
	sumSpans := 0
	for _, s := range rowSpans {
		fmt.Printf("Span: start %d, end %d\n", s.start, s.end)
		sumSpans += 1 + s.end - s.start
	}

	fmt.Printf("Result: spans = %d, objects = %v\n", sumSpans, objects)
}

func mDistance(aX, aY, bX, bY int) int {
	dX, dY := 0, 0
	if aX < bX {
		dX = bX - aX
	} else {
		dX = aX - bX
	}
	if aY < bY {
		dY = bY - aY
	} else {
		dY = aY - bY
	}
	return dX + dY
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
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
