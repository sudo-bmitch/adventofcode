package main

import (
	"bytes"
	"testing"
)

var day04Data = []byte(`MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`)

func TestDay04a(t *testing.T) {
	expect := "18"
	r := bytes.NewReader(day04Data)
	result, err := day04a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay04b(t *testing.T) {
	expect := "9"
	r := bytes.NewReader(day04Data)
	result, err := day04b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
