package main

import (
	"bytes"
	"testing"
)

var day05Data = []byte(`3-5
10-14
16-20
12-18

1
5
8
11
17
32`)

func TestDay05a(t *testing.T) {
	expect := "3"
	r := bytes.NewReader(day05Data)
	result, err := day05a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay05b(t *testing.T) {
	expect := "14"
	r := bytes.NewReader(day05Data)
	result, err := day05b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
