package main

import (
	"bytes"
	"testing"
)

var day08Data = []byte(`............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............`)

func TestDay08a(t *testing.T) {
	expect := "14"
	r := bytes.NewReader(day08Data)
	result, err := day08a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay08b(t *testing.T) {
	expect := "34"
	r := bytes.NewReader(day08Data)
	result, err := day08b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
