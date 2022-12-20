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

const (
	minutes = 24
)

type item int

const (
	itemNothing item = iota
	itemOre
	itemClay
	itemObsidian
	itemGeode
	itemCount
)

var debugLines = 0
var itemNames = [itemCount]string{"none", "ore", "clay", "obsidian", "geode"}

type state struct {
	robots    [itemCount]int
	inventory [itemCount]int
	build     item
	valid     [itemCount]bool
}

type blueprint struct {
	num  int
	cost [itemCount][itemCount]int
}

func main() {
	blueprints := []blueprint{}
	in := bufio.NewScanner(os.Stdin)
	reLine := regexp.MustCompile(`^Blueprint (\d+):\s+Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.$`)
	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		match := reLine.FindStringSubmatch(line)
		if len(match) <= 7 {
			fmt.Fprintf(os.Stderr, "failed parsing line: %s\n", line)
			return
		}
		bp := blueprint{}
		bp.num = mustAtoi(match[1])
		bp.cost[itemOre][itemOre] = mustAtoi(match[2])
		bp.cost[itemClay][itemOre] = mustAtoi(match[3])
		bp.cost[itemObsidian][itemOre] = mustAtoi(match[4])
		bp.cost[itemObsidian][itemClay] = mustAtoi(match[5])
		bp.cost[itemGeode][itemOre] = mustAtoi(match[6])
		bp.cost[itemGeode][itemObsidian] = mustAtoi(match[7])
		blueprints = append(blueprints, bp)
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	initState := state{}
	initState.robots[itemOre] = 1
	// check each plan for it's quality
	qualitySum := 0
	for planNum, plan := range blueprints {
		fmt.Printf("Reviewing blueprint %d\n", planNum)
		// initial state is one ore robot mining and no actions
		curState := [minutes]state{}
		for min := 0; min < minutes; min++ {
			if min == 0 {
				curState[min] = initState.nextState(&plan)
			} else {
				curState[min] = curState[min-1].nextState(&plan)
			}
		}
		fmt.Printf("initial state: \n%s\n", sprintStates(curState[:]))
		// bestState := [minutes]state{}
		// copy(bestState[:], curState[:])
		bestOutput := curState[minutes-1].inventory[itemGeode]
		for {
			// iterate last state, or try a previous state, until failing to iterate 0
			min := minutes - 1
			for min >= 0 {
				if curState[min].iterate(&plan) {
					break
				}
				min--
			}
			if min < 0 {
				break // finished iterating
			}
			// reset all following states
			for minNext := min + 1; minNext < minutes; minNext++ {
				curState[minNext] = curState[minNext-1].nextState(&plan)
			}
			// check if we have a new best state
			if curState[minutes-1].inventory[itemGeode] > bestOutput {
				bestOutput = curState[minutes-1].inventory[itemGeode]
				// fmt.Printf("found better: %d\n%s\n", bestOutput, sprintStates(curState[:]))
				fmt.Printf("found better: %d\n", bestOutput)
				// copy(bestState[:], curState[:]) // any need to track this?
			} else if debugLines > 0 {
				debugLines--
				fmt.Printf("debug\n%s\n", sprintStates(curState[:]))
			}
		}
		fmt.Printf("Plan %d has best output of %d.\n", planNum, bestOutput)
		// TODO: print state
		qualitySum += plan.num * bestOutput
	}

	fmt.Printf("Result: %d\n", qualitySum)
}

func (s *state) iterate(bp *blueprint) bool {
	if s.build == itemCount-1 {
		return false
	}
	// return inventory from item previous build item
	for i := s.build + 1; i < itemCount; i++ {
		if s.valid[i] {
			for j := range s.inventory {
				s.inventory[j] += bp.cost[s.build][j] - bp.cost[i][j]
			}
			s.build = i
			return true
		}
	}
	return false
}

func (s state) nextState(bp *blueprint) state {
	next := state{}
	copy(next.robots[:], s.robots[:])
	if s.build != itemNothing {
		next.robots[s.build]++
	}
	copy(next.inventory[:], s.inventory[:])
	// build first valid item
	next.valid = next.validBuild(bp)
	for i := itemNothing; i < itemCount; i++ {
		if next.valid[i] {
			next.build = i
			break
		}
	}
	for i := range next.inventory {
		next.inventory[i] += next.robots[i] - bp.cost[next.build][i]
	}
	return next
}

func (s state) validBuild(bp *blueprint) [itemCount]bool {
	valid := [itemCount]bool{}
	for i := itemNothing + 1; i < itemCount; i++ {
		able := true       // check if all materials are available
		enough := true     // check if no product needs more of a material
		generating := true // check if at least one robot exists for each material
		for j := itemNothing + 1; j < itemCount; j++ {
			if s.inventory[j] < bp.cost[i][j] {
				able = false // don't have inventory available to create
			}
			if s.inventory[i] < bp.cost[j][i] {
				enough = false // some product requires more of this
			}
			if bp.cost[j][i] > 0 && s.robots[i] == 0 {
				generating = false // not currently generating a needed item
			}
		}
		// can never have enough geodes :)
		if able {
			valid[i] = true
		} else if !enough && generating {
			valid[itemNothing] = true // wait is valid if we need more of something being generated
		}
	}
	return valid
}

func (s state) sprint() string {
	ret := fmt.Sprintf("build %s", itemNames[s.build])
	for i := 1; i < int(itemCount); i++ {
		ret += fmt.Sprintf(" %d/%d", s.robots[i], s.inventory[i])
	}
	return ret
}

func sprintStates(sl []state) string {
	ret := ""
	for min, s := range sl {
		ret += fmt.Sprintf("  min %d [%s]\n", min+1, s.sprint())
	}
	return ret
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
