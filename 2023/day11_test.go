package main

import (
	"bytes"
	"testing"
)

var day11Data = []byte(`...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....`)

func TestDay11a(t *testing.T) {
	expect := "374"
	r := bytes.NewReader(day11Data)
	result, err := day11a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay11b(t *testing.T) {
	tt := []struct {
		name   string
		expand int
		expect int
	}{
		{
			name:   "2",
			expand: 2,
			expect: 374,
		},
		{
			name:   "10",
			expand: 10,
			expect: 1030,
		},
		{
			name:   "100",
			expand: 100,
			expect: 8410,
		},
	}
	r := bytes.NewReader(day11Data)
	m, err := day11Parse(r)
	if err != nil {
		t.Errorf("failed to parse input: %v", err)
		return
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := day11Calc(m, tc.expand)
			if result != tc.expect {
				t.Errorf("unexpected result: expected %d, received %d", tc.expect, result)
			}
		})
	}
}
