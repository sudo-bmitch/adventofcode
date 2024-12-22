package main

import (
	"bytes"
	"testing"
)

var day22aData = []byte(`1
10
100
2024`)

var day22bData = []byte(`1
2
3
2024`)

func TestDay22a(t *testing.T) {
	expect := "37327623"
	r := bytes.NewReader(day22aData)
	result, err := day22a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay22b(t *testing.T) {
	expect := "-2,1,-1,3"
	r := bytes.NewReader(day22bData)
	result, err := day22b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
