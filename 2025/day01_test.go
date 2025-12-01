package main

import (
	"bytes"
	"testing"
)

var day01Data = []byte(`L68
L30
R48
L5
R60
L55
L1
L99
R14
L82`)

func TestDay01a(t *testing.T) {
	expect := "3"
	r := bytes.NewReader(day01Data)
	result, err := day01a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay01b(t *testing.T) {
	expect := "6"
	r := bytes.NewReader(day01Data)
	result, err := day01b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
