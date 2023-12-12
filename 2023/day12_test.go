package main

import (
	"bytes"
	"testing"
)

var day12Data = []byte(`???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1`)

func TestDay12a(t *testing.T) {
	expect := "21"
	r := bytes.NewReader(day12Data)
	result, err := day12a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay12b(t *testing.T) {
	expect := "525152"
	r := bytes.NewReader(day12Data)
	result, err := day12b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
