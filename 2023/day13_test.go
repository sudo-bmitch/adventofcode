package main

import (
	"bytes"
	"testing"
)

var day13Data = []byte(`#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.

#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#`)

func TestDay13a(t *testing.T) {
	expect := "405"
	r := bytes.NewReader(day13Data)
	result, err := day13a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay13b(t *testing.T) {
	expect := "400"
	r := bytes.NewReader(day13Data)
	result, err := day13b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
