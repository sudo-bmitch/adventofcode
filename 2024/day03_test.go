package main

import (
	"bytes"
	"testing"
)

var day03aData = []byte(`xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))`)

func TestDay03a(t *testing.T) {
	expect := "161"
	r := bytes.NewReader(day03aData)
	result, err := day03a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

var day03bData = []byte(`xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`)

func TestDay03b(t *testing.T) {
	expect := "48"
	r := bytes.NewReader(day03bData)
	result, err := day03b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
