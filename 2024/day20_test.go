package main

import (
	"bytes"
	"fmt"
	"testing"
)

var day20Data = []byte(`###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############`)

func TestDay20a(t *testing.T) {
	tt := []struct {
		minSave int
		expect  int
	}{
		{minSave: 65, expect: 0},
		{minSave: 64, expect: 1},
		{minSave: 40, expect: 2},
		{minSave: 38, expect: 3},
		{minSave: 36, expect: 4},
		{minSave: 20, expect: 5},
		{minSave: 12, expect: 8},
	}
	for _, tc := range tt {
		name := fmt.Sprintf("save-%d", tc.minSave)
		t.Run(name, func(t *testing.T) {
			r := bytes.NewReader(day20Data)
			result, err := day20Run(r, tc.minSave, 2)
			if err != nil {
				t.Errorf("failed with error: %v", err)
			} else if result != tc.expect {
				t.Errorf("unexpected result: expected %d, received %d", tc.expect, result)
			}

		})
	}
}

func TestDay20b(t *testing.T) {
	tt := []struct {
		minSave int
		expect  int
	}{
		{minSave: 77, expect: 0},
		{minSave: 76, expect: 3},
		{minSave: 75, expect: 3},
		{minSave: 74, expect: 7},
		{minSave: 73, expect: 7},
		{minSave: 72, expect: 29},
		{minSave: 71, expect: 29},
		{minSave: 70, expect: 41},
		{minSave: 69, expect: 41},
		{minSave: 68, expect: 55},
		{minSave: 67, expect: 55},
		{minSave: 66, expect: 67},
		{minSave: 65, expect: 67},
		{minSave: 64, expect: 86},
	}
	for _, tc := range tt {
		name := fmt.Sprintf("save-%d", tc.minSave)
		t.Run(name, func(t *testing.T) {
			r := bytes.NewReader(day20Data)
			result, err := day20Run(r, tc.minSave, 20)
			if err != nil {
				t.Errorf("failed with error: %v", err)
			} else if result != tc.expect {
				t.Errorf("unexpected result: expected %d, received %d", tc.expect, result)
			}

		})
	}
}
