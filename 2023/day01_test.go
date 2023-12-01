package main

import (
	"bytes"
	"testing"
)

var day01AData = []byte(`1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet`)

func TestDay01a(t *testing.T) {
	expect := "142"
	r := bytes.NewReader(day01AData)
	result, err := day01a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

var day01BData = []byte(`two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen`)

func TestDay01b(t *testing.T) {
	expect := "281"
	r := bytes.NewReader(day01BData)
	result, err := day01b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
