package main

import (
	"bytes"
	"testing"
)

var day14Data = []byte(`p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3`)

// var day14Data = []byte(`p=2,4 v=2,-3`)

func TestDay14a(t *testing.T) {
	expect := "12"
	r := bytes.NewReader(day14Data)
	result, err := day14a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay14b(t *testing.T) {
	t.Log("no tests available for 14b")
}
