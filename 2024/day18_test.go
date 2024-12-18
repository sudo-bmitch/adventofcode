package main

import (
	"bytes"
	"testing"
)

var day18Data = []byte(`5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0`)

func TestDay18a(t *testing.T) {
	expect := 22
	r := bytes.NewReader(day18Data)
	result, err := day18Run(7, 7, r, 12)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %d, received %d", expect, result)
	}
}

func TestDay18b(t *testing.T) {
	expect := "6,1"
	r := bytes.NewReader(day18Data)
	result, err := day18FindLimit(7, 7, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
