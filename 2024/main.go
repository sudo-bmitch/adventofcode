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

type dayFn func([]string, io.Reader) (string, error)

var days = map[string]dayFn{}

func registerDay(s string, fn dayFn) {
	days[s] = fn
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
