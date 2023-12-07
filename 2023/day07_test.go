package main

import (
	"bytes"
	"testing"
)

var day07Data = []byte(`32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483`)

func TestDay07a(t *testing.T) {
	expect := "6440"
	r := bytes.NewReader(day07Data)
	result, err := day07a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay07b(t *testing.T) {
	expect := "5905"
	r := bytes.NewReader(day07Data)
	result, err := day07b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
