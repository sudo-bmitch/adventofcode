package main

import (
	"bytes"
	"testing"
)

var day17Data = []byte(`2413432311323
3215453535623
3255245654254
3446585845452
4546657867536
1438598798454
4457876987766
3637877979653
4654967986887
4564679986453
1224686865563
2546548887735
4322674655533`)

func TestDay17a(t *testing.T) {
	expect := "102"
	r := bytes.NewReader(day17Data)
	result, err := day17a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay17b(t *testing.T) {
	expect := "94"
	r := bytes.NewReader(day17Data)
	result, err := day17b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
