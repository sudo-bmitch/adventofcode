package main

import (
	"bytes"
	"testing"
)

var day07Data = []byte(`.......S.......
...............
.......^.......
...............
......^.^......
...............
.....^.^.^.....
...............
....^.^...^....
...............
...^.^...^.^...
...............
..^...^.....^..
...............
.^.^.^.^.^...^.
...............`)

func TestDay07a(t *testing.T) {
	expect := "21"
	r := bytes.NewReader(day07Data)
	result, err := day07a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay07b(t *testing.T) {
	expect := "40"
	r := bytes.NewReader(day07Data)
	result, err := day07b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
