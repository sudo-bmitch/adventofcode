package main

import (
	"bytes"
	"testing"
)

var day12Data = []byte(`RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`)

func TestDay12a(t *testing.T) {
	expect := "1930"
	r := bytes.NewReader(day12Data)
	result, err := day12a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay12b(t *testing.T) {
	expect := "1206"
	r := bytes.NewReader(day12Data)
	result, err := day12b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
