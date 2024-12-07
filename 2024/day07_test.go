package main

import (
	"bytes"
	"testing"
)

var day07Data = []byte(`190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`)

func TestDay07a(t *testing.T) {
	expect := "3749"
	r := bytes.NewReader(day07Data)
	result, err := day07a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay07b(t *testing.T) {
	expect := "11387"
	r := bytes.NewReader(day07Data)
	result, err := day07b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
