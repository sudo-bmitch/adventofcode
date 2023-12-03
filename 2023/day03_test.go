package main

import (
	"bytes"
	"testing"
)

var data03Data = []byte(`467..114..
...*......
..35..633.
......#...
617*......
.....+.58.
..592.....
......755.
...$.*....
.664.598..`)

func TestDay03a(t *testing.T) {
	expect := "4361"
	r := bytes.NewReader(data03Data)
	result, err := day03a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay03b(t *testing.T) {
	expect := "467835"
	r := bytes.NewReader(data03Data)
	result, err := day03b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
