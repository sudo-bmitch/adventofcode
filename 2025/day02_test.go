package main

import (
	"bytes"
	"testing"
)

var day02Data = []byte(`11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124`)

func TestDay02a(t *testing.T) {
	expect := "1227775554"
	r := bytes.NewReader(day02Data)
	result, err := day02a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay02b(t *testing.T) {
	expect := "4174379265"
	r := bytes.NewReader(day02Data)
	result, err := day02b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
