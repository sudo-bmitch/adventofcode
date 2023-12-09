package main

import (
	"bufio"
	"fmt"
	"io"
)

func day09a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		nList, err := day09Parse(line)
		if err != nil {
			return "", err
		}
		dx := day09BuildDx(nList)
		// extrapolate
		for di := len(dx) - 1; di >= 0; di-- {
			if di >= len(dx)-1 {
				dx[di] = append(dx[di], 0)
			} else {
				dj := len(dx[di]) - 1
				dx[di] = append(dx[di], dx[di][dj]+dx[di+1][dj])
			}
		}
		sum += dx[0][len(dx[0])-1]
	}
	return fmt.Sprintf("%d", sum), nil
}

func day09b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		nList, err := day09Parse(line)
		if err != nil {
			return "", err
		}
		dx := day09BuildDx(nList)
		// extrapolate
		val := 0
		for di := len(dx) - 2; di >= 0; di-- {
			val = dx[di][0] - val
		}
		sum += val
	}
	return fmt.Sprintf("%d", sum), nil
}

func day09Parse(line string) ([]int, error) {
	return parseNumListBySpace(line)
}

func day09BuildDx(nList []int) [][]int {
	dx := [][]int{nList}
	di := 1
	for {
		dx = append(dx, []int{})
		allZero := true
		for prev := 0; prev < len(dx[di-1])-1; prev++ {
			entry := dx[di-1][prev+1] - dx[di-1][prev]
			dx[di] = append(dx[di], entry)
			if entry != 0 {
				allZero = false
			}
		}
		if allZero {
			break
		}
		di++
	}
	return dx
}
