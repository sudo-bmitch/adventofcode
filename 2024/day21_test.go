package main

import (
	"bytes"
	"testing"
)

var day21Data = []byte(`029A
980A
179A
456A
379A`)

func TestDay21a(t *testing.T) {
	expect := "126384"
	r := bytes.NewReader(day21Data)
	result, err := day21a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay21b(t *testing.T) {
	expect := "154115708116294"
	r := bytes.NewReader(day21Data)
	result, err := day21b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
