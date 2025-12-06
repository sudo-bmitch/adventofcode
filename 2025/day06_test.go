package main

import (
	"bytes"
	"testing"
)

var day06Data = []byte(`123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  `)

func TestDay06a(t *testing.T) {
	expect := "4277556"
	r := bytes.NewReader(day06Data)
	result, err := day06a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay06b(t *testing.T) {
	expect := "3263827"
	r := bytes.NewReader(day06Data)
	result, err := day06b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
