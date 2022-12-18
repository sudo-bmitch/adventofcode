package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	totalTime = 30
)

type valve struct {
	rate int
}
type tunnel struct {
	valves []string
}
type path struct {
	prev       *path  // previous steps back to nil at the start
	curName    string // name of valve
	arriveTime int    // time that you arrive at this location and finish opening the valve
	totalFlow  int    // sum of this and all previous open valves for the remainder
}

func main() {
	valves := map[string]valve{}
	tunnels := map[string]tunnel{}
	in := bufio.NewScanner(os.Stdin)
	lineRE := regexp.MustCompile(`^Valve ([A-Z]+) has flow rate=([0-9]+); tunnels? leads? to valves? (.*)$`)

	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		match := lineRE.FindStringSubmatch(line)
		if len(match) < 4 {
			fmt.Fprintf(os.Stderr, "line parsing failed for: %s\n", line)
			return
		}
		v := match[1]
		rate := mustAtoi(match[2])
		tunList := strings.Split(match[3], ", ")
		valves[v] = valve{rate: rate}
		tunnels[v] = tunnel{valves: tunList}
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	// get a list of all valves where you'd start a path
	startNames := []string{"AA"}
	for name, v := range valves {
		if v.rate > 0 {
			startNames = append(startNames, name)
		}
	}
	// get distances from starting points to any valve with a rate > 0
	distances := map[string]map[string]int{}
	for _, start := range startNames {
		distances[start] = map[string]int{}
		dFound := map[string]int{}
		curSearch := []string{start}
		for len(curSearch) > 0 {
			nextSearch := []string{}
			for _, cur := range curSearch {
				for _, next := range tunnels[cur].valves {
					if _, ok := dFound[next]; ok {
						continue
					}
					dFound[next] = dFound[cur] + 1
					if valves[next].rate > 0 {
						distances[start][next] = dFound[next]
					}
					nextSearch = append(nextSearch, next)
				}
			}
			curSearch = nextSearch
		}
	}

	paths := []*path{
		{
			curName: "AA",
		},
	}
	endPaths := []*path{}
	// loop through each minute
	for t := 1; t < totalTime; t++ {
		fmt.Printf("Searching at time %d, paths=%d\n", t, len(paths))
		nextPaths := []*path{}
		// loop through each path
		for _, p := range paths {
			p := p // can't remember if the pointer somehow gets reused
			// postpone paths we haven't arrived at yet
			if p.arriveTime > t {
				nextPaths = append(nextPaths, p)
				continue
			}
			// drop paths that aren't as good as the best, but only check after we reach their time
			if !p.bestPath() {
				continue
			}
			foundValid := false
			// try each route to each useful valve
			for dest, distance := range distances[p.curName] {
				if t+distance+1 >= totalTime {
					continue // path can't be reached in time
				}
				tryPath := path{
					prev:       p,
					curName:    dest,
					arriveTime: t + distance + 1,
					totalFlow:  p.totalFlow + valves[dest].rate*(totalTime-t-distance),
				}
				if tryPath.validPath() {
					nextPaths = append(nextPaths, &tryPath)
					foundValid = true
				}
			}
			if !foundValid && p.validPath() {
				endPaths = append(endPaths, p)
			}
		}
		paths = nextPaths
	}
	endPaths = append(endPaths, paths...)

	// look through paths to find best flow, show path
	bestI := 0
	bestRate := 0
	for i, p := range endPaths {
		if p.totalFlow > bestRate {
			bestI = i
			bestRate = p.totalFlow
		}
	}
	fmt.Printf("Best path: ")
	endPaths[bestI].print()
	fmt.Printf("Result: %d\n", bestRate)
}

func (p path) print() {
	cur := &p
	txt := []string{}
	for cur != nil {
		txt = append(txt, fmt.Sprintf("%s:%d ", cur.curName, cur.totalFlow))
		cur = cur.prev
	}
	for i := len(txt) - 1; i >= 0; i-- {
		fmt.Printf("%s", txt[i])
		if i > 0 {
			fmt.Printf(" ")
		}
	}
	fmt.Printf("\n")
}

var cachePaths = map[string]int{}

// return false if the path isn't valid
func (p path) validPath() bool {
	cur := p.prev
	for cur != nil {
		if cur.curName == p.curName {
			return false // valve already opened
		}
		cur = cur.prev
	}
	return true
}

func (p path) bestPath() bool {
	// check another path reached same endpoint with same open valves and a better flow rate
	valves := []string{}
	cur := &p
	for cur != nil {
		valves = append(valves, cur.curName)
		cur = cur.prev
	}
	sort.Strings(valves)
	key := fmt.Sprintf("%s:%s", p.curName, strings.Join(valves, ","))
	if cachePaths[key] > p.totalFlow {
		return false
	}
	cachePaths[key] = p.totalFlow
	return true
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
