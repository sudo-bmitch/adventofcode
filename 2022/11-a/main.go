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

const worryDiv = 3
const roundLen = 20

type monkey struct {
	items    []int64
	op       func(int64) int64
	test     func(int64) int
	inspects int
}

func monkeyNew(items []int64, opStr string, testDiv int, testTrue int, testFalse int) monkey {
	opParts := strings.Split(opStr, " ")
	if len(opParts) != 3 {
		panic(fmt.Sprintf("unhandled opStr: %s", opStr))
	}
	var arg1, arg2 int64
	if opParts[0] != "old" {
		i, err := strconv.Atoi(opParts[0])
		if err != nil {
			panic(err)
		}
		arg1 = int64(i)
	}
	if opParts[2] != "old" {
		i, err := strconv.Atoi(opParts[2])
		if err != nil {
			panic(err)
		}
		arg2 = int64(i)
	}
	opFn := func(old int64) int64 {
		if opParts[0] == "old" {
			arg1 = old
		}
		if opParts[2] == "old" {
			arg2 = old
		}
		switch opParts[1] {
		case "*":
			return arg1 * arg2
		case "+":
			return arg1 + arg2
		default:
			panic(fmt.Sprintf("unknown operation: %s", opParts[1]))
		}
	}
	testFn := func(old int64) int {
		if (old % int64(testDiv)) == 0 {
			return testTrue
		}
		return testFalse
	}
	return monkey{
		items: items,
		op:    opFn,
		test:  testFn,
	}
}

func main() {
	monkeys := []monkey{}

	reMonkey := regexp.MustCompile(`^Monkey (\d+):`)
	reItems := regexp.MustCompile("^Starting items: (.*)$")
	reOp := regexp.MustCompile("^Operation: new = (.*)$")
	reTest := regexp.MustCompile(`^Test: divisible by (\d+)$`)
	reTestTrue := regexp.MustCompile(`^If true: throw to monkey (\d+)$`)
	reTestFalse := regexp.MustCompile(`^If false: throw to monkey (\d+)$`)
	in := bufio.NewScanner(os.Stdin)
	// parse input
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		// parse monkey number
		monkeyMatch := reMonkey.FindStringSubmatch(line)
		if len(monkeyMatch) != 2 {
			fmt.Fprintf(os.Stderr, "Failed to parse input for monkey at: %s", line)
			return
		}
		monkeyNum, err := strconv.Atoi(monkeyMatch[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse input for monkey number at: %s", line)
			return
		}
		if monkeyNum != len(monkeys) {
			fmt.Fprintf(os.Stderr, "Out of order monkey, currently at %d: %s", len(monkeys), line)
			return
		}
		// parse starting items
		if !in.Scan() {
			break
		}
		line = strings.TrimSpace(in.Text())
		itemsMatch := reItems.FindStringSubmatch(line)
		if len(itemsMatch) != 2 {
			fmt.Fprintf(os.Stderr, "Failed to parse input for items at: %s", line)
			return
		}
		itemsSplit := strings.Split(itemsMatch[1], ", ")
		items := make([]int64, len(itemsSplit))
		for i, str := range itemsSplit {
			num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse input for item number %d at: %s: %v", i, line, err)
				return
			}
			items[i] = int64(num)
		}
		// parse operation
		if !in.Scan() {
			break
		}
		line = strings.TrimSpace(in.Text())
		opMatch := reOp.FindStringSubmatch(line)
		if len(opMatch) != 2 {
			fmt.Fprintf(os.Stderr, "Failed to parse input for operation at: %s", line)
			return
		}
		// parse test
		if !in.Scan() {
			break
		}
		line = strings.TrimSpace(in.Text())
		testMatch := reTest.FindStringSubmatch(line)
		if len(testMatch) != 2 {
			fmt.Fprintf(os.Stderr, "Failed to parse input for test at: %s", line)
			return
		}
		testDiv, err := strconv.Atoi(testMatch[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse input for test number at: %s", line)
			return
		}
		if !in.Scan() {
			break
		}
		line = strings.TrimSpace(in.Text())
		testTrueMatch := reTestTrue.FindStringSubmatch(line)
		if len(testTrueMatch) != 2 {
			fmt.Fprintf(os.Stderr, "Failed to parse input for testTrue at: %s", line)
			return
		}
		testTrue, err := strconv.Atoi(testTrueMatch[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse input for testTrue number at: %s", line)
			return
		}
		if !in.Scan() {
			break
		}
		line = strings.TrimSpace(in.Text())
		testFalseMatch := reTestFalse.FindStringSubmatch(line)
		if len(testFalseMatch) != 2 {
			fmt.Fprintf(os.Stderr, "Failed to parse input for testFalse at: %s", line)
			return
		}
		testFalse, err := strconv.Atoi(testFalseMatch[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse input for testFalse number at: %s", line)
			return
		}
		// create monkey and append to list
		monkeys = append(monkeys, monkeyNew(items, opMatch[1], testDiv, testTrue, testFalse))
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	// perform the rounds
	for round := 0; round < roundLen; round++ {
		for monkeyI, monkey := range monkeys {
			for _, item := range monkey.items {
				worry := int64(monkey.op(item) / worryDiv)
				target := monkey.test(worry)
				if target < 0 || target >= len(monkeys) || target == monkeyI {
					panic(fmt.Sprintf("out of bounds throw target: from %d, to %d, worry %d", monkeyI, target, worry))
				}
				monkeys[target].items = append(monkeys[target].items, worry)
			}
			monkeys[monkeyI].inspects += len(monkey.items)
			monkeys[monkeyI].items = []int64{} // all items thrown
		}
		// report result
		fmt.Printf("After round %d:\n", round+1)
		for monkeyI, monkey := range monkeys {
			fmt.Printf("Monkey %d: ", monkeyI)
			for i, item := range monkey.items {
				fmt.Printf("%d", item)
				if i+1 != len(monkey.items) {
					fmt.Printf(", ")
				}
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}

	// get most active two monkeys
	inspects := []int{}
	for monkeyI, monkey := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times.\n", monkeyI, monkey.inspects)
		inspects = append(inspects, monkey.inspects)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(inspects)))

	fmt.Printf("Result: %d\n", inspects[0]*inspects[1])
}
