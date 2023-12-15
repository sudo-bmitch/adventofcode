package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	err := runDay()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v", err)
		os.Exit(1)
	}
}

func runDay() error {
	day := ""
	args := []string{}
	if len(os.Args) > 1 {
		day = os.Args[1]
		args = os.Args[2:]
	} else {
		fmt.Printf("Enter a day: ")
		_, err := fmt.Scanf("%s", &day)
		if err != nil {
			return err
		}
	}
	if fn, ok := days[day]; ok {
		result, err := fn(args, os.Stdin)
		if err != nil {
			return err
		}
		fmt.Printf("result: %s\n", result)
	} else {
		return fmt.Errorf("day not found: %s", day)
	}

	return nil
}

var days = map[string]func([]string, io.Reader) (string, error){
	"01a": day01a,
	"01b": day01b,
	"02a": day02a,
	"02b": day02b,
	"03a": day03a,
	"03b": day03b,
	"04a": day04a,
	"04b": day04b,
	"05a": day05a,
	"05b": day05b,
	"06a": day06a,
	"06b": day06b,
	"07a": day07a,
	"07b": day07b,
	"08a": day08a,
	"08b": day08b,
	"09a": day09a,
	"09b": day09b,
	"10a": day10a,
	"10b": day10b,
	"11a": day11a,
	"11b": day11b,
	"12a": day12a,
	"12b": day12b,
	"13a": day13a,
	"13b": day13b,
	"14a": day14a,
	"14b": day14b,
	"15a": day15a,
	"15b": day15b,
}
