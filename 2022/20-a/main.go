package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var debug = false

type packet struct {
	value int
	moved bool
}

func main() {
	packets := []packet{}
	in := bufio.NewScanner(os.Stdin)
	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		val := mustAtoi(line)
		packets = append(packets, packet{value: val})
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	// run packet mix
	numPackets := len(packets)
	for i := 0; i < numPackets; i++ {
		if packets[i].moved {
			continue
		}
		if packets[i].value%numPackets == 0 {
			// noop
			packets[i].moved = true
		}
		// remove from list and reinsert in new offset
		p := packets[i]
		p.moved = true
		packets = append(packets[:i], packets[i+1:]...)
		newI := i + p.value
		for newI < 0 {
			newI += (numPackets - 1)
		}
		newI = newI % (numPackets - 1)
		// circular list can add to beginning or end, and end is easier
		if newI == 0 || newI >= len(packets) {
			packets = append(packets, p)
		} else {
			// add a space to the list
			packets = append(packets[:newI+1], packets[newI:]...)
			packets[newI] = p
		}
		if debug {
			fmt.Printf("Debug: after moving %d to %d = ", p.value, newI)
			printPackets(packets)
		}
		// inserted after present position, so next i should be unchanged
		if newI > i {
			i--
		}
	}

	if debug {
		fmt.Printf("Debug: list = ")
		printPackets(packets)
	}

	// find the 0 value
	zeroI := 0
	for i, p := range packets {
		if p.value == 0 {
			zeroI = i
			break
		}
	}
	if debug {
		fmt.Printf("Debug: zero position = %d\n", zeroI)
	}

	result := 0
	result += packets[(zeroI+1000)%numPackets].value
	result += packets[(zeroI+2000)%numPackets].value
	result += packets[(zeroI+3000)%numPackets].value
	fmt.Printf("Result: %d\n", result)
}

func printPackets(packets []packet) {
	numPackets := len(packets)
	for i, p := range packets {
		fmt.Printf("%d", p.value)
		if i != numPackets-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("\n")
}

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
