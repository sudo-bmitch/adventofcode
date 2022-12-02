package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// first col: opponent A=rock, B=paper, C=scissors
// second col: you X=rock, Y=paper, Z=scissors
// lose=0, draw=3, win=6, rock=1, paper=2, scissors=3
var pointsLU = map[string]int{
	"A X": 3 + 1,
	"A Y": 6 + 2,
	"A Z": 0 + 3,
	"B X": 0 + 1,
	"B Y": 3 + 2,
	"B Z": 6 + 3,
	"C X": 6 + 1,
	"C Y": 0 + 2,
	"C Z": 3 + 3,
}

func main() {
	total := 0
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		lu, ok := pointsLU[line]
		if !ok {
			fmt.Fprintf(os.Stderr, "Failed to lookup line: %s\n", line)
		}
		total += lu
	}
	fmt.Printf("Total: %d\n", total)

}
