package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type day06Race struct {
	t, d int
}

func day06a(args []string, rdr io.Reader) (string, error) {
	races, err := day06aParse(rdr)
	if err != nil {
		return "", err
	}
	perm := 1 // permutations of winning races
	for _, r := range races {
		span, err := day06CalcSpan(r)
		if err != nil {
			return "", err
		}
		perm *= span
	}
	return fmt.Sprintf("%d", perm), nil
}

func day06b(args []string, rdr io.Reader) (string, error) {
	races, err := day06bParse(rdr) // different parser for second part
	if err != nil {
		return "", err
	}
	perm := 1 // permutations of winning races
	for _, r := range races {
		span, err := day06CalcSpan(r)
		if err != nil {
			return "", err
		}
		perm *= span
	}
	return fmt.Sprintf("%d", perm), nil
}

func day06aParse(rdr io.Reader) ([]day06Race, error) {
	in := bufio.NewScanner(rdr)
	if !in.Scan() {
		return nil, fmt.Errorf("failed to scan times")
	}
	line := in.Text()
	if line[:6] != "Time: " {
		return nil, fmt.Errorf("time line is missing prefix: %s", line)
	}
	times, err := parseNumListBySpace(line[6:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse times from %s: %w", line, err)
	}

	if !in.Scan() {
		return nil, fmt.Errorf("failed to scan distances")
	}
	line = in.Text()
	if line[:10] != "Distance: " {
		return nil, fmt.Errorf("distance line is missing prefix: %s", line)
	}
	distances, err := parseNumListBySpace(line[10:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse distances from %s: %w", line, err)
	}

	if len(times) != len(distances) {
		return nil, fmt.Errorf("time and distance counts do not match: %v, %v", times, distances)
	}
	races := []day06Race{}
	for i := range times {
		races = append(races, day06Race{t: times[i], d: distances[i]})
	}

	return races, nil
}

func day06bParse(rdr io.Reader) ([]day06Race, error) {
	in := bufio.NewScanner(rdr)
	if !in.Scan() {
		return nil, fmt.Errorf("failed to scan times")
	}
	line := in.Text()
	if line[:6] != "Time: " {
		return nil, fmt.Errorf("time line is missing prefix: %s", line)
	}
	t, err := strconv.Atoi(strings.ReplaceAll(line[6:], " ", ""))
	if err != nil {
		return nil, fmt.Errorf("failed to parse time from %s: %w", line, err)
	}

	if !in.Scan() {
		return nil, fmt.Errorf("failed to scan distances")
	}
	line = in.Text()
	if line[:10] != "Distance: " {
		return nil, fmt.Errorf("distance line is missing prefix: %s", line)
	}
	d, err := strconv.Atoi(strings.ReplaceAll(line[10:], " ", ""))
	if err != nil {
		return nil, fmt.Errorf("failed to parse distance from %s: %w", line, err)
	}

	return []day06Race{{t: t, d: d}}, nil
}

func day06CalcDist(r day06Race, hold int) int {
	if hold >= r.t || hold < 0 {
		return 0
	}
	return (r.t - hold) * hold
}

func day06CalcSpan(r day06Race) (int, error) {
	// binary search for any entry above r.d
	l, m, h := 0, r.t/2, r.t
	for {
		dCur := day06CalcDist(r, m)
		if dCur > r.d {
			break
		}
		dPrev := day06CalcDist(r, m-1)
		if dPrev < dCur {
			l = m + 1
		} else if dPrev > dCur {
			h = m - 1
		} else {
			return 1, fmt.Errorf("binary search unable to pick next direction, l=%d, m=%d, h=%d, dCur=%d", l, m, h, dCur)
		}
		if l >= h {
			return 1, fmt.Errorf("binary search for faster race failed: l=%d, m=%d, h=%d, dCur=%d", l, m, h, dCur)
		}
		m = l + ((h - l) / 2)
	}
	// save M to reset for the limHi search
	saveM := m
	// find the lower limit
	l, h = 0, m
	for {
		m = l + ((h - l) / 2) // default to round down when h-l = 1
		dCur := day06CalcDist(r, m)
		if dCur > r.d {
			h = m
		} else {
			l = m + 1
		}
		if l >= h {
			break
		}
	}
	limLo := h
	// find the upper limit
	l, h = saveM, r.t
	for {
		m = l + ((h - l + 1) / 2) // +1 to round up when h-l = 1
		dCur := day06CalcDist(r, m)
		if dCur > r.d {
			l = m
		} else {
			h = m - 1
		}
		if l >= h {
			break
		}
	}
	limHi := l
	if limHi-limLo+1 < 1 {
		return 1, fmt.Errorf("failed to find a valid span, race time %d, goal %d, limLo=%d, limHi=%d", r.t, r.d, limLo, limHi)
	}
	return limHi - limLo + 1, nil
}
