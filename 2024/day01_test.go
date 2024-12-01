package main

import (
	"bytes"
	"testing"
)

var day01Data = []byte(`3   4
4   3
2   5
1   3
3   9
3   3`)

func TestDay01a(t *testing.T) {
	expect := "11"
	r := bytes.NewReader(day01Data)
	result, err := day01a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay01b(t *testing.T) {
	expect := "31"
	r := bytes.NewReader(day01Data)
	result, err := day01b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
