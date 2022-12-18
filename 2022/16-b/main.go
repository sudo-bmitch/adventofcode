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
	totalTime = 26
	searchers = 2
)

type valve struct {
	rate int
}
type tunnel struct {
	valves []string
}
type step struct {
	name      string
	index     int
	time      int
	totalRate int
}
type path []step
type pathSet [searchers]path

var (
	valves    = map[string]valve{}
	distances = map[string]map[string]int{}
)

func main() {
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

	curPathSet := pathSet{{{name: "AA", time: 1}}, {{name: "AA", time: 1}}}
	bestPathSet := pathSet{}
	copy(bestPathSet[0], curPathSet[0])
	copy(bestPathSet[1], curPathSet[1])
	bestRate := 0
	valid := [searchers]bool{true, true}
	listValves := []string{}
	for name, v := range valves {
		if v.rate > 0 {
			listValves = append(listValves, name)
		}
	}
	sort.Strings(listValves)

	for {
		// increment curPathSet
		// if 0 isn't valid, reset v1, increment last element and check valid, else pop and leave invalid
		// else if 1 is empty, append next unused item to 0 and check valid
		// else if 1 isn't valid, increment last element and check valid, else pop and leave invalid
		// else append next unused item to 1 and check valid
		// to increment, find index of last item in list, go through remaining until unused found

		if len(curPathSet[0]) == 1 && !valid[0] {
			// exhausted all items to search
			break
		} else if !valid[0] {
			// when first searcher is invalid or finished search on 1
			// reset 1 when changing 0
			curPathSet[1] = path{{name: "AA", time: 1}}
			valid[1] = true
			// find the next element for end of list
			pos := len(curPathSet[0]) - 1
			i, ok := curPathSet.unusedNext(curPathSet[0][pos].index+1, listValves)
			if ok {
				name := listValves[i]
				start := curPathSet[0][pos-1].time + distances[curPathSet[0][pos-1].name][name]
				curPathSet[0][pos] = step{
					name:      name,
					index:     i,
					time:      start + 1,
					totalRate: curPathSet[0][pos-1].totalRate + valves[name].rate*(totalTime-start),
				}
				valid[0] = curPathSet[0].validPath()
			} else {
				// else no more options for last element, pop, and leave invalid to increment previous position
				curPathSet[0] = curPathSet[0][:len(curPathSet[0])-1]
			}
		} else if len(curPathSet[1]) == 1 && !valid[1] {
			// if second searcher has exhausted their search, append to first searcher
			// reset 1 when changing 0
			curPathSet[1] = path{{name: "AA", time: 1}}
			valid[1] = true
			// append next unused item to 0
			i, ok := curPathSet.unusedNext(0, listValves)
			if ok {
				pos := len(curPathSet[0])
				name := listValves[i]
				start := curPathSet[0][pos-1].time + distances[curPathSet[0][pos-1].name][name]
				curPathSet[0] = append(curPathSet[0], step{
					name:      name,
					index:     i,
					time:      start + 1,
					totalRate: curPathSet[0][pos-1].totalRate + valves[name].rate*(totalTime-start),
				})
				valid[0] = curPathSet[0].validPath()
			} else {
				// no more unused locations, mark invalid to increment/pop
				valid[0] = false
			}
		} else if !valid[1] {
			// increment last element
			pos := len(curPathSet[1]) - 1
			i, ok := curPathSet.unusedNext(curPathSet[1][pos].index+1, listValves)
			if ok {
				name := listValves[i]
				start := curPathSet[1][pos-1].time + distances[curPathSet[1][pos-1].name][name]
				curPathSet[1][pos] = step{
					name:      name,
					index:     i,
					time:      start + 1,
					totalRate: curPathSet[1][pos-1].totalRate + valves[name].rate*(totalTime-start),
				}
				valid[1] = curPathSet[1].validPath()
			} else {
				// else pop and leave invalid
				curPathSet[1] = curPathSet[1][:len(curPathSet[1])-1]
			}
		} else {
			// append next unused item to 1
			i, ok := curPathSet.unusedNext(0, listValves)
			if ok {
				pos := len(curPathSet[1])
				name := listValves[i]
				start := curPathSet[1][pos-1].time + distances[curPathSet[1][pos-1].name][name]
				curPathSet[1] = append(curPathSet[1], step{
					name:      name,
					index:     i,
					time:      start + 1,
					totalRate: curPathSet[1][pos-1].totalRate + valves[name].rate*(totalTime-start),
				})
				valid[1] = curPathSet[1].validPath()
			} else {
				// no more unused locations, mark invalid to increment/pop
				valid[1] = false
			}
		}

		// check if better than best
		if valid[0] && valid[1] {
			if bestRate < curPathSet.total() {
				bestRate = curPathSet.total()
				fmt.Printf("New best: ")
				curPathSet.print()
				bestPathSet[0] = make(path, len(curPathSet[0]))
				bestPathSet[1] = make(path, len(curPathSet[1]))
				copy(bestPathSet[0], curPathSet[0])
				copy(bestPathSet[1], curPathSet[1])
			}
		}
	}

	fmt.Printf("Best path: ")
	bestPathSet.print()
	fmt.Printf("Result: %d\n", bestPathSet.total())
}

func (ps pathSet) print() {
	fmt.Printf("total=%d: ", ps.total())
	for i, p := range ps {
		for j, e := range p {
			fmt.Printf("%s", e.name)
			if j != len(p)-1 {
				fmt.Printf("->")
			}
		}
		if i != len(ps)-1 {
			fmt.Printf(" / ")
		} else {
			fmt.Printf("\n")
		}
	}
}

func (ps pathSet) total() int {
	return ps[0][len(ps[0])-1].totalRate + ps[1][len(ps[1])-1].totalRate
}

func (ps pathSet) unusedNext(start int, list []string) (int, bool) {
	for i := start; i < len(list); i++ {
		if !ps.used(list[i]) {
			return i, true
		}
	}
	return -1, false
}

func (ps pathSet) used(name string) bool {
	for _, p := range ps {
		for _, e := range p {
			if e.name == name {
				return true
			}
		}
	}
	return false
}

// return false if the path isn't valid
func (p path) validPath() bool {
	// ensure time doesn't exceed limit
	if len(p) == 0 {
		return false
	}
	if p[len(p)-1].time > totalTime {
		return false
	}
	return true
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
