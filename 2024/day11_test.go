package main

import (
	"bytes"
	"testing"
)

var day11Data = []byte(`125 17`)

func TestDay11a(t *testing.T) {
	expect := "55312"
	r := bytes.NewReader(day11Data)
	result, err := day11a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay11b(t *testing.T) {
	expect := "65601038650482"
	r := bytes.NewReader(day11Data)
	result, err := day11b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
