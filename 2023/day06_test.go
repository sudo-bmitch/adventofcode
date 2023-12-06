package main

import (
	"bytes"
	"testing"
)

var day06Data = []byte(`Time:      7  15   30
Distance:  9  40  200`)

func TestDay06a(t *testing.T) {
	expect := "288"
	r := bytes.NewReader(day06Data)
	result, err := day06a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay06b(t *testing.T) {
	expect := "71503"
	r := bytes.NewReader(day06Data)
	result, err := day06b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
