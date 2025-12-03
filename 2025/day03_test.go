package main

import (
	"bytes"
	"testing"
)

var day03Data = []byte(`987654321111111
811111111111119
234234234234278
818181911112111`)

func TestDay03a(t *testing.T) {
	expect := "357"
	r := bytes.NewReader(day03Data)
	result, err := day03a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay03b(t *testing.T) {
	expect := "3121910778619"
	r := bytes.NewReader(day03Data)
	result, err := day03b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
