package main

import (
	"bytes"
	"testing"
)

var day09Data = []byte(`7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`)

func TestDay09a(t *testing.T) {
	expect := "50"
	r := bytes.NewReader(day09Data)
	result, err := day09a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay09b(t *testing.T) {
	expect := "24"
	r := bytes.NewReader(day09Data)
	result, err := day09b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
