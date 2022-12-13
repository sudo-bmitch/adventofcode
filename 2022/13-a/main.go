package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type data struct {
	val  int
	list []*data
}

func cmpData(a, b data) int {
	// compare value
	if a.list == nil && b.list == nil {
		if a.val > b.val {
			return 1
		} else if a.val < b.val {
			return -1
		}
		return 0
	}
	// upgrade val to list
	if a.list == nil {
		a.list = []*data{{val: a.val}}
	}
	if b.list == nil {
		b.list = []*data{{val: b.val}}
	}
	// compare lists
	for i := range a.list {
		if i >= len(b.list) {
			return 1
		}
		if cmp := cmpData(*a.list[i], *b.list[i]); cmp != 0 {
			return cmp
		}
	}
	if len(a.list) < len(b.list) {
		return -1
	}
	return 0
}

type parseData struct {
	data     data
	valueSet bool
	parent   *parseData
}

func parsePacket(str string) data {
	state := "start"
	cur := &parseData{}
	for _, c := range str {
		switch state {
		case "start":
			switch c {
			case '[':
				state = "head"
			default:
				panic("parse error")
			}
		case "head": // start of list
			cur.data.list = []*data{}
			cur = &parseData{parent: cur}
			switch {
			case c == '[':
				// repeat head
			case c == ']':
				// empty list
				state = "tail"
			case c >= '0' && c <= '9':
				// value
				val := int(rune(c) - rune('0'))
				if val < 0 || val > 9 {
					panic("bad number parsing")
				}
				cur.data.val = val
				state = "value"
			default:
				panic("parse error")
			}
		case "value": // within a number
			cur.valueSet = true
			switch {
			case c == ']':
				state = "tail"
			case c == ',':
				// end value, start of next
				state = "sep"
			case c >= '0' && c <= '9':
				val := int(rune(c) - rune('0'))
				if val < 0 || val > 9 {
					panic("bad number parsing")
				}
				cur.data.val = cur.data.val*10 + val
			default:
				panic("unknown value")
			}
		case "sep": // separator (comma)
			// add cur data to the parent list and create a new data
			cur.parent.data.list = append(cur.parent.data.list, &cur.data)
			cur = &parseData{parent: cur.parent}
			switch {
			case c == '[':
				state = "head"
			case c >= '0' && c <= '9':
				// value
				val := int(rune(c) - rune('0'))
				if val < 0 || val > 9 {
					panic("bad number parsing")
				}
				cur.data.val = val
				state = "value"
			default:
				panic("parse error")
			}
		case "tail": // end of list
			if cur.parent == nil {
				panic("too many closing brackets")
			}
			// end list
			if cur.valueSet || cur.data.list != nil {
				cur.parent.data.list = append(cur.parent.data.list, &cur.data)
			}
			cur = cur.parent
			switch {
			case c == ',':
				state = "sep"
			case c == ']':
				// repeat tail
			default:
				panic("parse error")
			}
		default:
			panic("unknown state")
		}
	}
	// handle last tail entry
	if state != "tail" || cur.parent == nil {
		panic("parse error")
	}
	if cur.valueSet || cur.data.list != nil {
		cur.parent.data.list = append(cur.parent.data.list, &cur.data)
	}
	cur = cur.parent
	return cur.data
}

func printData(p data) {
	if p.list != nil {
		fmt.Printf("[")
		for i, e := range p.list {
			printData(*e)
			if i < len(p.list)-1 {
				fmt.Printf(",")
			}
		}
		fmt.Printf("]")
	} else {
		fmt.Printf("%d", p.val)
	}
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	index := 1
	sum := 0
	// parse input
	for {
		if !in.Scan() {
			break
		}
		line1 := strings.TrimSpace(in.Text())
		if line1 == "" {
			continue
		}
		if !in.Scan() {
			break
		}
		line2 := strings.TrimSpace(in.Text())
		if line2 == "" {
			fmt.Fprintf(os.Stderr, "unexpected empty line after %s\n", line1)
			return
		}
		data1 := parsePacket(line1)
		data2 := parsePacket(line2)
		fmt.Printf("data1: ")
		printData(data1)
		fmt.Printf("\n")
		fmt.Printf("data2: ")
		printData(data2)
		fmt.Printf("\n")
		if cmpData(data1, data2) < 0 {
			fmt.Printf("pair %d in right order\n", index)
			sum += index
		}
		index++
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}
	fmt.Printf("Result: %d\n", sum)
}
