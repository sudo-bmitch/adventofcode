package main

import (
	"bytes"
	"testing"
)

var day10Data = []byte(`89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`)

func TestDay10a(t *testing.T) {
	expect := "36"
	r := bytes.NewReader(day10Data)
	result, err := day10a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay10b(t *testing.T) {
	expect := "81"
	r := bytes.NewReader(day10Data)
	result, err := day10b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
