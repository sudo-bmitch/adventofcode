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

const (
	decryptKey = 811589153
	mixCount   = 10
)

type packet struct {
	value      int
	prev, next *packet
}

func main() {
	// packets is a double linked list
	packets := []*packet{}
	in := bufio.NewScanner(os.Stdin)
	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		val := mustAtoi(line)
		p := &packet{value: val * decryptKey}
		if len(packets) > 0 {
			p.next = packets[0]
			p.prev = packets[len(packets)-1]
			packets[0].prev = p
			packets[len(packets)-1].next = p
		} else {
			p.next = p
			p.prev = p
		}
		packets = append(packets, p)
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	numPackets := len(packets)
	if debug {
		fmt.Printf("input packets = ")
		printPackets(packets)
	}
	// run packet mix
	for mc := 0; mc < mixCount; mc++ {
		for i := 0; i < numPackets; i++ {
			if packets[i].value%numPackets == 0 {
				// noop
				continue
			}
			// remove from list
			p := packets[i]
			prev := p.prev
			prev.next = p.next
			p.next.prev = prev
			// get count to shift and walk prev pointer
			shift := p.value % (numPackets - 1)
			if shift > 0 {
				for n := 0; n < shift; n++ {
					prev = prev.next
				}
			} else {
				for n := 0; n > shift; n-- {
					prev = prev.prev
				}
			}
			// reinsert in new position
			p.prev = prev
			p.next = prev.next
			p.next.prev = p
			p.prev.next = p
		}
		if debug {
			fmt.Printf("Debug: packets after %d rounds = ", mc)
			printPackets(packets)
		} else {
			fmt.Printf("Finished round %d\n", mc)
		}
	}

	// find the 0 value
	var zeroP *packet
	for _, p := range packets {
		if p.value == 0 {
			zeroP = p
			break
		}
	}
	if zeroP == nil {
		fmt.Fprintf(os.Stderr, "failed to find zero value\n")
		return
	}
	result := 0
	// walk next pointers by 1000%len
	shift := 1000 % numPackets
	cur := zeroP
	for _, n := range []int{1, 2, 3} {
		for count := 0; count < shift; count++ {
			cur = cur.next
		}
		if debug {
			fmt.Printf("Result value %d is %d\n", n, cur.value)
		}
		result += cur.value
	}
	fmt.Printf("Result: %d\n", result)
}

func printPackets(packets []*packet) {
	// print by walking head until returning to head
	if len(packets) > 0 {
		head := packets[0]
		cur := head
		for {
			fmt.Printf("%d", cur.value)
			cur = cur.next
			if cur == head {
				break
			}
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
