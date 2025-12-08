package main

import (
	"bytes"
	"testing"
)

var day08Data = []byte(`162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689`)

func TestDay08a(t *testing.T) {
	expect := "40"
	day08aLimit = 10
	r := bytes.NewReader(day08Data)
	result, err := day08a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay08b(t *testing.T) {
	expect := "25272"
	r := bytes.NewReader(day08Data)
	result, err := day08b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
