package main

import (
	"bytes"
	"testing"
)

var dayXXData = []byte(`1234`)

func TestDayXXa(t *testing.T) {
	expect := "42"
	r := bytes.NewReader(dayXXData)
	result, err := dayXXa([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDayXXb(t *testing.T) {
	expect := "42"
	r := bytes.NewReader(dayXXData)
	result, err := dayXXb([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
