package main

import (
	"bytes"
	"testing"
)

var day09Data = []byte(`0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45`)

func TestDay09a(t *testing.T) {
	expect := "114"
	r := bytes.NewReader(day09Data)
	result, err := day09a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay09b(t *testing.T) {
	expect := "2"
	r := bytes.NewReader(day09Data)
	result, err := day09b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
