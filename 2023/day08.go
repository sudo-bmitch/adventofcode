package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type day08Map struct {
	instruct []rune
	nodes    map[string]day08LRNode
}

type day08LRNode struct {
	l, r string
}

type day08Path struct {
	startup  int
	zList    []int
	zI       int
	cycleLen int
	cycleI   int
}

func day08a(args []string, rdr io.Reader) (string, error) {
	step := 0
	m, err := day08Parse(rdr)
	if err != nil {
		return "", err
	}
	cur := "AAA"
	for cur != "ZZZ" {
		n, ok := m.nodes[cur]
		if !ok {
			return "", fmt.Errorf("failed to lookup node for %s on step %d", cur, step)
		}
		switch m.instruct[step%len(m.instruct)] {
		case 'L':
			cur = n.l
		case 'R':
			cur = n.r
		default:
			return "", fmt.Errorf("unknown instruct direction %c", m.instruct[step%len(m.instruct)])
		}
		step++
	}
	return fmt.Sprintf("%d", step), nil
}

func day08b(args []string, rdr io.Reader) (string, error) {
	m, err := day08Parse(rdr)
	if err != nil {
		return "", err
	}
	startList := []string{}
	for k := range m.nodes {
		if k[2] == 'A' {
			startList = append(startList, k)
		}
	}
	// treat each start as a wheel
	// every hit on a Z entry is tracked as a cycle point as offset from first Z entry
	// when loop detected with interval x*len(path), switch to tracking next point in cycle for that start
	// when all starts have a loop, switch to skipping each wheel < the highest forward one offset point
	// stop when all wheels == same highest point
	paths := make([]day08Path, len(startList))
	for i, start := range startList {
		paths[i], err = day08BuildPath(m, start)
		if err != nil {
			return "", err
		}
	}
	curSteps := make([]int, len(startList))
	lastStep := 1 // force all paths to init on first pass (curSteps[*] = 0)
	for {
		atEnd := 0
		for i := range paths {
			if curSteps[i] < lastStep {
				curSteps[i] = paths[i].Next()
				if curSteps[i] > lastStep {
					lastStep = curSteps[i]
					atEnd = 0
				}
			}
			if curSteps[i] == lastStep {
				atEnd++
			}
		}
		if atEnd == len(paths) {
			break
		}
	}
	return fmt.Sprintf("%d", lastStep), nil
}

var day08NodeRE = regexp.MustCompile(`^([A-Z1-3]+) = \(([A-Z1-3]+), ([A-Z1-3]+)\)$`)

func day08Parse(rdr io.Reader) (day08Map, error) {
	m := day08Map{
		nodes: map[string]day08LRNode{},
	}
	in := bufio.NewScanner(rdr)
	if !in.Scan() {
		return m, fmt.Errorf("missing instructions")
	}
	line := strings.TrimSpace(in.Text())
	m.instruct = []rune(line)
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		match := day08NodeRE.FindStringSubmatch(line)
		if len(match) != 4 {
			return m, fmt.Errorf("node match on %s failed, len=%d", line, len(match))
		}
		m.nodes[match[1]] = day08LRNode{l: match[2], r: match[3]}
	}
	return m, nil
}

func day08BuildPath(m day08Map, start string) (day08Path, error) {
	p := day08Path{}
	cur := start
	step := 0
	for cur[2] != 'Z' {
		n, ok := m.nodes[cur]
		if !ok {
			return p, fmt.Errorf("failed to lookup node for %s on instance %s step %d", cur, start, step)
		}
		switch m.instruct[step%len(m.instruct)] {
		case 'L':
			cur = n.l
		case 'R':
			cur = n.r
		default:
			return p, fmt.Errorf("unknown instruct direction %c", m.instruct[step%len(m.instruct)])
		}
		step++
	}
	fmt.Fprintf(os.Stderr, "startup for %s at %d\n", start, step)
	p.startup = step
	p.zList = []int{0}
	for {
		n, ok := m.nodes[cur]
		if !ok {
			return p, fmt.Errorf("failed to lookup node for %s on instance %s step %d", cur, start, step)
		}
		switch m.instruct[step%len(m.instruct)] {
		case 'L':
			cur = n.l
		case 'R':
			cur = n.r
		default:
			return p, fmt.Errorf("unknown instruct direction %c", m.instruct[step%len(m.instruct)])
		}
		step++
		if cur[2] == 'Z' {
			if (step-p.startup)%len(m.instruct) == 0 {
				fmt.Fprintf(os.Stderr, "cycle length for %s on %s is %d (mod %d)\n", start, cur, step-p.startup, (step-p.startup)%len(m.instruct))
				p.cycleLen = step - p.startup
				break // found the end of the loop
			}
			fmt.Fprintf(os.Stderr, "add zList for %s on %s at %d (mod %d)\n", start, cur, step-p.startup, (step-p.startup)%len(m.instruct))
			p.zList = append(p.zList, step-p.startup)
		}
	}
	return p, nil
}

func (p *day08Path) Next() int {
	n := p.startup + (p.cycleLen * p.cycleI) + p.zList[p.zI]
	p.zI++
	if p.zI >= len(p.zList) {
		p.zI = 0
		p.cycleI++
	}
	return n
}
