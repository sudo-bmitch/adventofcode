package main

import (
	"bytes"
	"testing"
)

var day19Data = []byte(`r, wr, b, g, bwu, rb, gb, br

brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb`)

func TestDay19a(t *testing.T) {
	expect := "6"
	r := bytes.NewReader(day19Data)
	result, err := day19a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay19b(t *testing.T) {
	expect := "16"
	r := bytes.NewReader(day19Data)
	result, err := day19b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
