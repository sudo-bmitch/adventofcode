package main

import (
	"fmt"
	"io"
)

func init() {
	registerDay("XXa", dayXXa)
	registerDay("XXb", dayXXb)
}

func dayXXa(args []string, rdr io.Reader) (string, error) {
	sum := 0
	// read list of numbers
	// for split := range parse.MustNumSlice(rdr) {
	// }

	// read a grid of runes
	// g, err := grid.FromReader(rdr)
	// if err != nil {
	// 	return "", err
	// }

	// read lines individually
	// in := bufio.NewScanner(rdr)
	// for in.Scan() {
	// 	// parse each line into a solution and list of numbers
	// 	line := in.Text()
	// }

	return fmt.Sprintf("%d", sum), nil
}

func dayXXb(args []string, rdr io.Reader) (string, error) {
	sum := 0

	return fmt.Sprintf("%d", sum), nil
}
