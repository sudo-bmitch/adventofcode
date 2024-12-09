package main

import (
	"bytes"
	"testing"
)

var day09Data = []byte(`2333133121414131402`)

func TestDay09a(t *testing.T) {
	expect := "1928"
	r := bytes.NewReader(day09Data)
	result, err := day09a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay09b(t *testing.T) {
	expect := "2858"
	r := bytes.NewReader(day09Data)
	result, err := day09b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
