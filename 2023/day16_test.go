package main

import (
	"bytes"
	"testing"
)

var day16Data = []byte(`.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....`)

func TestDay16a(t *testing.T) {
	expect := "46"
	r := bytes.NewReader(day16Data)
	result, err := day16a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay16b(t *testing.T) {
	expect := "51"
	r := bytes.NewReader(day16Data)
	result, err := day16b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
