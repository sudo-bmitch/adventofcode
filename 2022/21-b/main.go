package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var debug = true

type monkey struct {
	name     string
	set      bool
	value    int
	inA, inB string
	op       operation
}
type operation int

const (
	opAdd operation = iota
	opSub
	opMul
	opDiv
	opEq
)

var opLU = map[byte]operation{
	'+': opAdd,
	'-': opSub,
	'*': opMul,
	'/': opDiv,
	'=': opEq,
}
var opStr = map[operation]string{
	opAdd: "+",
	opSub: "-",
	opMul: "*",
	opDiv: "/",
	opEq:  "=",
}

func main() {
	monkeys := map[string]*monkey{}
	in := bufio.NewScanner(os.Stdin)
	lineRE := regexp.MustCompile(`^([a-z]*): (([0-9]+)|(([a-z]+) ([\-\+\*\/\=]) ([a-z]+)))$`)
	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		match := lineRE.FindStringSubmatch(line)
		if len(match) < 8 {
			fmt.Fprintf(os.Stderr, "failed to parse %s: %v\n", line, match)
		}
		m := &monkey{name: match[1]}
		if len(match[3]) > 0 {
			m.value = mustAtoi(match[3])
			m.set = true
		}
		if len(match[4]) > 0 {
			m.inA, m.inB = match[5], match[7]
			m.op = opLU[match[6][0]]
		}
		monkeys[match[1]] = m
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	pending := map[string][]string{}
	ready := []string{}
	for name := range monkeys {
		set, dep := monkeys[name].trySet(monkeys)
		if set {
			if debug {
				fmt.Printf("debug: set: %s\n", monkeys[name].string())
			}
			if nowReady, ok := pending[name]; ok {
				ready = append(ready, nowReady...)
				delete(pending, name)
			}
		} else {
			if debug {
				fmt.Printf("debug: dep %s: %s\n", dep, monkeys[name].string())
			}
			pending[dep] = append(pending[dep], name)
		}
	}
	for len(ready) > 0 {
		name := ready[0]
		ready = ready[1:]
		set, dep := monkeys[name].trySet(monkeys)
		if set {
			if debug {
				fmt.Printf("debug: set: %s\n", monkeys[name].string())
			}
			if nowReady, ok := pending[name]; ok {
				ready = append(ready, nowReady...)
				delete(pending, name)
			}
		} else {
			if debug {
				fmt.Printf("debug: dep %s: %s\n", dep, monkeys[name].string())
			}
			pending[dep] = append(pending[dep], name)
		}
	}

	// solve for humn
	nextName := "root"
	nextVal := 0
	can := true
	for nextName != "humn" && can {
		if debug {
			fmt.Printf("debug: can set %s to %d: ", nextName, nextVal)
		}
		can, nextName, nextVal = monkeys[nextName].canSet(nextVal, monkeys)
		if debug {
			fmt.Printf("%t if %s is %d\n", can, nextName, nextVal)
		}
	}
	if !can || nextName != "humn" {
		fmt.Printf("stuck trying to find human value\n")
		return
	}

	fmt.Printf("Result: %d\n", nextVal)
}

func (m monkey) string() string {
	if m.set {
		return fmt.Sprintf("%s: %d", m.name, m.value)
	} else {
		return fmt.Sprintf("%s: %s %s %s", m.name, m.inA, opStr[m.op], m.inB)
	}
}

// trySet attempts to set the value for a monkey
func (m *monkey) trySet(monkeys map[string]*monkey) (bool, string) {
	if m.set {
		return true, ""
	}
	if monkeys[m.inA] == nil || !monkeys[m.inA].set {
		return false, m.inA
	} else if monkeys[m.inB] == nil || !monkeys[m.inB].set {
		return false, m.inB
	}
	switch m.op {
	case opAdd:
		m.value = monkeys[m.inA].value + monkeys[m.inB].value
	case opSub:
		m.value = monkeys[m.inA].value - monkeys[m.inB].value
	case opMul:
		m.value = monkeys[m.inA].value * monkeys[m.inB].value
	case opDiv:
		m.value = monkeys[m.inA].value / monkeys[m.inB].value
	}
	m.set = true
	return true, ""
}

// canSet indicates what monkey needs to be set to a value
func (m *monkey) canSet(value int, monkeys map[string]*monkey) (bool, string, int) {
	if m.set {
		return value == m.value, "", 0
	}
	switch m.op {
	case opEq:
		if monkeys[m.inA] != nil && monkeys[m.inA].set {
			if monkeys[m.inB] != nil && monkeys[m.inB].set {
				return monkeys[m.inA] == monkeys[m.inB], "", 0
			}
			return true, m.inB, monkeys[m.inA].value
		} else if monkeys[m.inB] != nil && monkeys[m.inB].set {
			return true, m.inA, monkeys[m.inB].value
		} else {
			return false, "", 0
		}
	case opAdd:
		if monkeys[m.inA] != nil && monkeys[m.inA].set {
			if monkeys[m.inB] != nil && monkeys[m.inB].set {
				return monkeys[m.inA].value+monkeys[m.inB].value == value, "", 0
			}
			return true, m.inB, value - monkeys[m.inA].value
		} else if monkeys[m.inB] != nil && monkeys[m.inB].set {
			return true, m.inA, value - monkeys[m.inB].value
		} else {
			return false, "", 0
		}
	case opSub:
		if monkeys[m.inA] != nil && monkeys[m.inA].set {
			if monkeys[m.inB] != nil && monkeys[m.inB].set {
				return monkeys[m.inA].value-monkeys[m.inB].value == value, "", 0
			}
			return true, m.inB, monkeys[m.inA].value - value
		} else if monkeys[m.inB] != nil && monkeys[m.inB].set {
			return true, m.inA, value + monkeys[m.inB].value
		} else {
			return false, "", 0
		}
	case opMul:
		if monkeys[m.inA] != nil && monkeys[m.inA].set {
			if monkeys[m.inB] != nil && monkeys[m.inB].set {
				return monkeys[m.inA].value*monkeys[m.inB].value == value, "", 0
			}
			return true, m.inB, value / monkeys[m.inA].value
		} else if monkeys[m.inB] != nil && monkeys[m.inB].set {
			return true, m.inA, value / monkeys[m.inB].value
		} else {
			return false, "", 0
		}
	case opDiv:
		if monkeys[m.inA] != nil && monkeys[m.inA].set {
			if monkeys[m.inB] != nil && monkeys[m.inB].set {
				return monkeys[m.inA].value/monkeys[m.inB].value == value, "", 0
			}
			return true, m.inB, monkeys[m.inA].value / value
		} else if monkeys[m.inB] != nil && monkeys[m.inB].set {
			return true, m.inA, value * monkeys[m.inB].value
		} else {
			return false, "", 0
		}
	}
	return false, "", 0
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
