package main

import (
	"bytes"
	"testing"
)

var day14Data = []byte(`O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....`)

func TestDay14a(t *testing.T) {
	expect := "136"
	r := bytes.NewReader(day14Data)
	result, err := day14a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay14b(t *testing.T) {
	expect := "64"
	r := bytes.NewReader(day14Data)
	result, err := day14b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
