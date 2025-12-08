package main

import (
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
)

func init() {
	registerDay("08a", day08a)
	registerDay("08b", day08b)
}

type day08Box struct {
	X, Y, Z int
}

var (
	day08aLimit = 1000
)

func day08a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	debug := false
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	boxes := []day08Box{}
	for line := range strings.SplitSeq(strings.TrimSpace(string(in)), "\n") {
		numStrList := strings.SplitN(line, ",", 3)
		if len(numStrList) != 3 {
			return "", fmt.Errorf("failed to split line %s", line)
		}
		x, err := strconv.Atoi(numStrList[0])
		if err != nil {
			return "", fmt.Errorf("failed to convert to number %s from line %s: %w", numStrList[0], line, err)
		}
		y, err := strconv.Atoi(numStrList[1])
		if err != nil {
			return "", fmt.Errorf("failed to convert to number %s from line %s: %w", numStrList[1], line, err)
		}
		z, err := strconv.Atoi(numStrList[2])
		if err != nil {
			return "", fmt.Errorf("failed to convert to number %s from line %s: %w", numStrList[2], line, err)
		}
		boxes = append(boxes, day08Box{X: x, Y: y, Z: z})
	}
	lenList := day08ConnLen(boxes)
	connections := map[int]map[int]bool{}
	count := 0
	// connect together the requested number of boxes
	for _, curLen := range slices.Sorted(maps.Keys(lenList)) {
		for _, connPair := range lenList[curLen] {
			a, b := connPair[0], connPair[1]
			// skip existing connections
			if connections[a] != nil && connections[a][b] {
				continue
			}
			// track connections made
			if debug {
				fmt.Fprintf(os.Stderr, "connecting %d (%v) with %d (%v), lenSq %d\n", a, boxes[a], b, boxes[b], curLen)
			}
			if connections[a] == nil {
				connections[a] = map[int]bool{}
			}
			if connections[b] == nil {
				connections[b] = map[int]bool{}
			}
			connections[a][b] = true
			connections[b][a] = true
			count++
			if count >= day08aLimit {
				break
			}
		}
		if count >= day08aLimit {
			break
		}
	}
	// track the number of connections per group
	groupSizes := []int{}
	seen := map[int]bool{}
	for a := range connections {
		if seen[a] {
			continue
		}
		size := 0
		queue := []int{a}
		for len(queue) > 0 {
			// pop the tail
			cur := queue[len(queue)-1]
			queue = queue[:len(queue)-1]
			if seen[cur] {
				continue
			}
			seen[cur] = true
			size++
			for b := range connections[cur] {
				queue = append(queue, b)
			}
		}
		groupSizes = append(groupSizes, size)
	}
	slices.Sort(groupSizes)
	slices.Reverse(groupSizes)
	if len(groupSizes) < 3 {
		return "", fmt.Errorf("not enough groups, only have the following sizes: %v", groupSizes)
	}
	sum = groupSizes[0] * groupSizes[1] * groupSizes[2]
	return fmt.Sprintf("%d", sum), nil
}

func day08b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	debug := false
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	boxes := []day08Box{}
	for line := range strings.SplitSeq(strings.TrimSpace(string(in)), "\n") {
		numStrList := strings.SplitN(line, ",", 3)
		if len(numStrList) != 3 {
			return "", fmt.Errorf("failed to split line %s", line)
		}
		x, err := strconv.Atoi(numStrList[0])
		if err != nil {
			return "", fmt.Errorf("failed to convert to number %s from line %s: %w", numStrList[0], line, err)
		}
		y, err := strconv.Atoi(numStrList[1])
		if err != nil {
			return "", fmt.Errorf("failed to convert to number %s from line %s: %w", numStrList[1], line, err)
		}
		z, err := strconv.Atoi(numStrList[2])
		if err != nil {
			return "", fmt.Errorf("failed to convert to number %s from line %s: %w", numStrList[2], line, err)
		}
		boxes = append(boxes, day08Box{X: x, Y: y, Z: z})
	}
	connections := map[int]map[int]bool{}
	lastA, lastB := 0, 0
	lenList := day08ConnLen(boxes)
	// initialize from the first node
	seen := map[int]bool{0: true}
	for _, curLen := range slices.Sorted(maps.Keys(lenList)) {
		for _, connPair := range lenList[curLen] {
			a, b := connPair[0], connPair[1]
			// skip existing connections
			if connections[a] != nil && connections[a][b] {
				continue
			}
			// track connections made
			if debug {
				fmt.Fprintf(os.Stderr, "connecting %d (%v) with %d (%v), lenSq %d\n", a, boxes[a], b, boxes[b], curLen)
			}
			if connections[a] == nil {
				connections[a] = map[int]bool{}
			}
			if connections[b] == nil {
				connections[b] = map[int]bool{}
			}
			connections[a][b] = true
			connections[b][a] = true
			// if either side is in the fully connected list, extend the seen map
			if seen[a] || seen[b] {
				queue := []int{a, b}
				for len(queue) > 0 {
					// pop the tail
					cur := queue[len(queue)-1]
					queue = queue[:len(queue)-1]
					// skip nodes we've previously seen
					if seen[cur] {
						continue
					}
					// extend via existing connections
					seen[cur] = true
					for next := range connections[cur] {
						queue = append(queue, next)
					}
				}
				if debug {
					fmt.Fprintf(os.Stderr, "seen list at %d nodes\n", len(seen))
				}
			}
			if len(seen) == len(boxes) {
				// we've seen every box, last connection was a to b
				lastA = boxes[a].X
				lastB = boxes[b].X
				break
			}
		}
		if len(seen) == len(boxes) {
			break
		}
	}
	sum = lastA * lastB
	return fmt.Sprintf("%d", sum), nil
	// 7264308110
}

func day08LenSq(a, b day08Box) int {
	return (a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y) + (a.Z-b.Z)*(a.Z-b.Z)
}

// create a map of connection lengths (squared), each map is a slice of connections between two boxes.
func day08ConnLen(boxes []day08Box) map[int][][2]int {
	lenList := map[int][][2]int{}
	for a := range boxes {
		for b := a + 1; b < len(boxes); b++ {
			lenSq := day08LenSq(boxes[a], boxes[b])
			lenList[lenSq] = append(lenList[lenSq], [2]int{a, b})
		}
	}
	return lenList
}
