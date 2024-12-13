package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	registerDay("13a", day13a)
	registerDay("13b", day13b)
}

// a * ax + b * bx = px
// a * ay + b * by = py

// b = (px - a * ax) / bx
// b = (py - a * ay) / by
// (px - a * ax) / bx = (py - a * ay) / by
// px * by - a * ax * by = py * bx - a * ay * bx
// a * ay * bx - a * ax * by = py * bx - px * by
// a * (ay * bx - ax * by) = py * bx - px * by
// a = (py * bx - px * by) / (ay * bx - ax * by)

// a = (px - b * bx) / ax
// a = (py - b * by) / ay
// (px - b * bx) / ax = (py - b * by) / ay
// (px - b * bx) * ay = (py - b * by) * ax
// px * ay - b * bx * ay = py * ax - b * by * ax
// b * by * ax - b * bx * ay = py * ax - px * ay
// b * (by * ax - bx * ay) = py * ax - px * ay
// b = (py * ax - px * ay) / (by * ax - bx * ay)

// Example 1 works
// ax=94, ay=34, bx=22, by=67, px=8400, py=5400
// a = (5400 * 22 - 8400 * 67) / (34 * 22 - 94 * 67) = 80
// b = (5400 * 94 - 8400 * 34) / (67 * 94 - 22 * 34) = 40

// Example 2, expected fail
// ax=26, ay=66, bx=67, by=21, px=12748, py=12176
// a = (12176 * 67 - 12748 * 21) / (66 * 67 - 26 * 21) = 141.404...
// b = (12176 * 26 - 12748 * 66) / (21 * 26 - 67 * 66) = 135.395...

func day13a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	costA := 3
	costB := 1
	limit := 100
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	games := strings.Split(string(in), "\n\n")
	gameRE := regexp.MustCompile(`^\s*Button A: X\+(\d+), Y\+(\d+)\s+Button B: X\+(\d+), Y\+(\d+)\s+Prize: X=(\d+), Y=(\d+)\s*$`)
	for _, game := range games {
		gameParse := gameRE.FindStringSubmatch(game)
		if len(gameParse) != 7 {
			return "", fmt.Errorf("game did not parse, %s: %v", game, gameParse)
		}
		ax, err := strconv.Atoi(gameParse[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		ay, err := strconv.Atoi(gameParse[2])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		bx, err := strconv.Atoi(gameParse[3])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		by, err := strconv.Atoi(gameParse[4])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		px, err := strconv.Atoi(gameParse[5])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		py, err := strconv.Atoi(gameParse[6])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		af := float64(py*bx-px*by) / float64(ay*bx-ax*by)
		bf := float64(py*ax-px*ay) / float64(by*ax-bx*ay)
		a := int(af + 0.5)
		b := int(bf + 0.5)
		if a > 0 && b > 0 && a <= limit && b <= limit && a*ax+b*bx == px && a*ay+b*by == py {
			sum += costA*a + costB*b
		}

	}

	return fmt.Sprintf("%d", sum), nil
}

func day13b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	costA := 3
	costB := 1
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	games := strings.Split(string(in), "\n\n")
	gameRE := regexp.MustCompile(`^\s*Button A: X\+(\d+), Y\+(\d+)\s+Button B: X\+(\d+), Y\+(\d+)\s+Prize: X=(\d+), Y=(\d+)\s*$`)
	for _, game := range games {
		gameParse := gameRE.FindStringSubmatch(game)
		if len(gameParse) != 7 {
			return "", fmt.Errorf("game did not parse, %s: %v", game, gameParse)
		}
		ax, err := strconv.Atoi(gameParse[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		ay, err := strconv.Atoi(gameParse[2])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		bx, err := strconv.Atoi(gameParse[3])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		by, err := strconv.Atoi(gameParse[4])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		px, err := strconv.Atoi(gameParse[5])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		py, err := strconv.Atoi(gameParse[6])
		if err != nil {
			return "", fmt.Errorf("failed to parse number %s: %w", gameParse[1], err)
		}
		px += 10000000000000
		py += 10000000000000
		af := float64(py*bx-px*by) / float64(ay*bx-ax*by)
		bf := float64(py*ax-px*ay) / float64(by*ax-bx*ay)
		a := int(af + 0.5)
		b := int(bf + 0.5)
		if a > 0 && b > 0 && a*ax+b*bx == px && a*ay+b*by == py {
			sum += costA*a + costB*b
		}

	}

	return fmt.Sprintf("%d", sum), nil
}
