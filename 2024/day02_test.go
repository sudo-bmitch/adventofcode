package main

import (
	"bytes"
	"testing"
)

var day02Data = []byte(`7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9`)

func TestDay02a(t *testing.T) {
	expect := "2"
	r := bytes.NewReader(day02Data)
	result, err := day02a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay02b(t *testing.T) {
	expect := "4"
	r := bytes.NewReader(day02Data)
	result, err := day02b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
