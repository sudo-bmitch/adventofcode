package main

import (
	"bytes"
	"testing"
)

var day15Data = []byte(`rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7`)

func TestDay15a(t *testing.T) {
	expect := "1320"
	r := bytes.NewReader(day15Data)
	result, err := day15a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay15b(t *testing.T) {
	expect := "145"
	r := bytes.NewReader(day15Data)
	result, err := day15b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
