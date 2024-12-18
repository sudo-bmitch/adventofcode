package main

import (
	"bytes"
	"testing"
)

var day17DataA = []byte(`Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0`)
var day17DataB = []byte(`Register A: 2024
Register B: 0
Register C: 0

Program: 0,3,5,4,3,0`)

func TestDay17a(t *testing.T) {
	expect := "4,6,3,5,6,3,5,2,1,0"
	r := bytes.NewReader(day17DataA)
	result, err := day17a([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
	tt := []struct {
		name             string
		regA, regB, regC int
		program          []int
		expect           string
		exA, exB, exC    int
	}{
		{
			name:    "ex1",
			regC:    9,
			program: []int{2, 6, 5, 5},
			expect:  "1",
		},
		{
			name:    "ex2",
			regA:    10,
			program: []int{5, 0, 5, 1, 5, 4},
			expect:  "0,1,2",
		},
		{
			name:    "ex3",
			regA:    2024,
			program: []int{0, 1, 5, 4, 3, 0, 5, 4},
			expect:  "4,2,5,6,7,7,7,7,3,1,0,0",
		},
		{
			name:    "ex4",
			regB:    29,
			program: []int{1, 7},
			exB:     26,
		},
		{
			name:    "ex5",
			regB:    2024,
			regC:    43690,
			program: []int{4, 0},
			exB:     44354,
		},
		{
			name:    "part2",
			regA:    117440,
			program: []int{0, 3, 5, 4, 3, 0},
			expect:  "0,3,5,4,3,0",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tVM := day17VM{
				regA:    tc.regA,
				regB:    tc.regB,
				regC:    tc.regC,
				program: tc.program,
			}
			err := tVM.Run()
			if err != nil {
				t.Fatalf("failed to run: %v", err)
			}
			out := tVM.PrintOut()
			if out != tc.expect {
				t.Errorf("output, expected %s, received %s", tc.expect, out)
			}
			if tc.exA != 0 && tc.exA != tVM.regA {
				t.Errorf("register A, expected %d, received %d", tc.exA, tVM.regA)
			}
			if tc.exB != 0 && tc.exB != tVM.regB {
				t.Errorf("register B, expected %d, received %d", tc.exB, tVM.regB)
			}
			if tc.exC != 0 && tc.exC != tVM.regC {
				t.Errorf("register C, expected %d, received %d", tc.exC, tVM.regC)
			}
		})
	}
}

func TestDay17b(t *testing.T) {
	expect := "117440"
	r := bytes.NewReader(day17DataB)
	result, err := day17b([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}

func TestDay17slow(t *testing.T) {
	expect := "117440"
	r := bytes.NewReader(day17DataB)
	result, err := day17slow([]string{}, r)
	if err != nil {
		t.Errorf("failed with error: %v", err)
	} else if result != expect {
		t.Errorf("unexpected result: expected %s, received %s", expect, result)
	}
}
