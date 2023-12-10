package main

import (
	"bytes"
	"testing"
)

var day10Data = []byte(`..F7.
.FJ|.
SJ.L7
|F--J
LJ...`)

func TestDay10a(t *testing.T) {
	expect := "8"
	r := bytes.NewReader(day10Data)
	result, err := day10a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay10b(t *testing.T) {
	tt := []struct {
		name   string
		data   []byte
		expect string
	}{
		{
			name: "one",
			data: []byte(`...........
.S-------7.
.|F-----7|.
.||.....||.
.||.....||.
.|L-7.F-J|.
.|..|.|..|.
.L--J.L--J.
...........`),
			expect: "4",
		},
		{
			name: "two",
			data: []byte(`.F----7F7F7F7F-7....
.|F--7||||||||FJ....
.||.FJ||||||||L7....
FJL7L7LJLJ||LJ.L-7..
L--J.L7...LJS7F-7L7.
....F-J..F7FJ|L7L7L7
....L7.F7||L7|.L7L7|
.....|FJLJ|FJ|F7|.LJ
....FJL-7.||.||||...
....L---J.LJ.LJLJ...`),
			expect: "8",
		},
		{
			name: "three",
			data: []byte(`FF7FSF7F7F7F7F7F---7
L|LJ||||||||||||F--J
FL-7LJLJ||||||LJL-77
F--JF--7||LJLJ7F7FJ-
L---JF-JLJ.||-FJLJJ7
|F|F-JF---7F7-L7L|7|
|FFJF7L7F-JF7|JL---7
7-L-JL7||F7|L7F-7F7|
L.L7LFJ|||||FJL7||LJ
L7JLJL-JLJLJL--JLJ.L`),
			expect: "10",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.data)
			result, err := day10b([]string{}, r)
			if err != nil {
				t.Errorf("failed with error: %v", err)
			} else if result != tc.expect {
				t.Errorf("unexpected result: expected %s, received %s", tc.expect, result)
			}
		})
	}
}
