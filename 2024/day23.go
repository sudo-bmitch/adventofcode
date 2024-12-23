package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func init() {
	registerDay("23a", day23a)
	registerDay("23b", day23b)
}

func day23a(args []string, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(in), "\n")
	nodeLinks := map[string][]string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		nodes := strings.Split(line, "-")
		if len(nodes) != 2 {
			return "", fmt.Errorf("unknown line: %s", line)
		}
		nodeLinks[nodes[0]] = append(nodeLinks[nodes[0]], nodes[1])
		nodeLinks[nodes[1]] = append(nodeLinks[nodes[1]], nodes[0])
	}

	interconnect := map[string]bool{}
	for node := range nodeLinks {
		// only search nodes starting with a t
		if node[0] != 't' {
			continue
		}
		for i, link1 := range nodeLinks[node] {
			for _, link2 := range nodeLinks[node][i+1:] {
				// check if link1 is linked to link2
				if day23IsLinked(nodeLinks, link1, link2) {
					interconnect[day23NodeList(node, link1, link2)] = true
				}
			}
		}
	}

	sum := len(interconnect)
	return fmt.Sprintf("%d", sum), nil
}

func day23b(args []string, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(in), "\n")
	nodeLinks := map[string][]string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		nodes := strings.Split(line, "-")
		if len(nodes) != 2 {
			return "", fmt.Errorf("unknown line: %s", line)
		}
		nodeLinks[nodes[0]] = append(nodeLinks[nodes[0]], nodes[1])
		nodeLinks[nodes[1]] = append(nodeLinks[nodes[1]], nodes[0])
	}

	bestList := []string{}
	for node := range nodeLinks {
		if len(bestList) >= len(nodeLinks[node])+1 {
			continue
		}
		// all possible groupings from this node
		first := append([]string{node}, nodeLinks[node]...)
		possibleGroups := [][]string{first}
	groupLoop: // not to be confused with a fruit loop
		for curList := 0; curList < len(possibleGroups); curList++ {
			// check if every node is connected to all remaining nodes in the list
			for node1 := 1; node1 < len(possibleGroups[curList]); node1++ {
				for node2 := node1 + 1; node2 < len(possibleGroups[curList]); node2++ {
					if !day23IsLinked(nodeLinks, possibleGroups[curList][node1], possibleGroups[curList][node2]) {
						// when two nodes aren't linked
						// stop if removing an entry would make this list too small
						if len(bestList) >= len(possibleGroups[curList])-1 {
							continue groupLoop
						}
						// else duplicate the possible group and remove one node from each
						dup := make([]string, len(possibleGroups[curList]))
						copy(dup, possibleGroups[curList])
						if node2 <= len(possibleGroups[curList])-1 {
							possibleGroups[curList] = possibleGroups[curList][:node2]
						} else {
							possibleGroups[curList] = append(possibleGroups[curList][:node2], possibleGroups[curList][node2+1:]...)
						}
						dup = append(dup[:node1], dup[node1+1:]...)
						possibleGroups = append(possibleGroups, dup)
					}
				}
			}
			// track the best
			if len(possibleGroups[curList]) > len(bestList) {
				bestList = possibleGroups[curList]
			}
		}
	}
	return day23NodeList(bestList...), nil
}

func day23IsLinked(nodeLinks map[string][]string, node1, node2 string) bool {
	for _, dest := range nodeLinks[node1] {
		if dest == node2 {
			return true
		}
	}
	return false
}

func day23NodeList(nodes ...string) string {
	sort.Strings(nodes)
	return strings.Join(nodes, ",")
}
