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

type xy struct {
	x, y int
}

func main() {
	knots := [10]xy{}
	// easier than handling the slice/append logic for each x-y entry, I may regret this later
	visited := map[string]bool{"0,0": true}

	in := bufio.NewScanner(os.Stdin)
	// read moves
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// split by direction and distance
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "invalid input (space): %s\n", line)
			return
		}
		dir := parts[0]
		steps, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid input (number): %s\n", line)
			return
		}
		diff := xy{}
		switch dir {
		case "R":
			diff.x = 1
		case "L":
			diff.x = -1
		case "U":
			diff.y = 1
		case "D":
			diff.y = -1
		}
		// perform move
		move(&knots, diff, &visited, steps)
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}
	count := len(visited)
	fmt.Printf("Result: %d\n", count)
}

func follow(a, b xy) xy {
	diff := xy{}
	dx := a.x - b.x
	dy := a.y - b.y
	// if dx or dy diff is 2, shift 1 in that direction
	if dx > 1 {
		diff.x = 1
	} else if dx < -1 {
		diff.x = -1
	}
	if dy > 1 {
		diff.y = 1
	} else if dy < -1 {
		diff.y = -1
	}
	// when shifting 1 in either direction, always move diagonal towards target
	if diff.y != 0 && diff.x == 0 && dx != 0 {
		diff.x = dx
	}
	if diff.x != 0 && diff.y == 0 && dy != 0 {
		diff.y = dy
	}
	return diff
}

func move(knots *[10]xy, diff xy, visited *map[string]bool, count int) {
	for i := 0; i < count; i++ {
		knots[0].x += diff.x
		knots[0].y += diff.y
		for k := 1; k < 10; k++ {
			d := follow(knots[k-1], knots[k])
			knots[k].x += d.x
			knots[k].y += d.y
		}
		pos := fmt.Sprintf("%d,%d", knots[9].x, knots[9].y)
		if !(*visited)[pos] {
			(*visited)[pos] = true
		}
	}
}
