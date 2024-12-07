package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func init() {
	registerDay("07a", day07a)
	registerDay("07b", day07b)
}

func day07a(args []string, rdr io.Reader) (string, error) {
	return day07Run(args, rdr, day07OpMul)
}
func day07b(args []string, rdr io.Reader) (string, error) {
	return day07Run(args, rdr, day07OpConcat)
}

func day07Run(_ []string, rdr io.Reader, lastOp day07Ops) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		// parse each line into a solution and list of numbers
		line := in.Text()
		splitStr := strings.Fields(line)
		if len(splitStr) == 0 {
			continue
		}
		if len(splitStr) < 2 || splitStr[0][len(splitStr[0])-1] != ':' {
			return "", fmt.Errorf("invalid input: %s", line)
		}
		solution, err := strconv.Atoi(splitStr[0][:len(splitStr[0])-1])
		if err != nil {
			return "", fmt.Errorf("failed to parse solution %s: %w", line, err)
		}
		numbers := make([]int, len(splitStr)-1)
		for i := range numbers {
			numbers[i], err = strconv.Atoi(splitStr[i+1])
			if err != nil {
				return "", fmt.Errorf("failed to parse number[%d] %s: %w", i, line, err)
			}
		}
		ops := make([]day07Ops, len(numbers)-1)
		// iterate over each combination of operations
		for {
			tally := numbers[0]
			for i, op := range ops {
				switch op {
				case day07OpAdd:
					tally += numbers[i+1]
				case day07OpMul:
					tally *= numbers[i+1]
				case day07OpConcat:
					max := 1
					for numbers[i+1] >= max {
						tally *= 10
						max *= 10
					}
					tally += numbers[i+1]
				default:
					return "", fmt.Errorf("unknown op")
				}
			}
			// check the result
			if tally == solution {
				sum += solution
				break
			}
			// switch to the next combination of ops
			if day07OpsNext(ops, lastOp) {
				break
			}
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

type day07Ops int

const (
	day07OpAdd = iota
	day07OpMul
	day07OpConcat
)

// day07OpsNext returns true when it is done iterating and has returned to all Adds.
func day07OpsNext(opList []day07Ops, opLast day07Ops) bool {
	carry := true
	for i := 0; i < len(opList) && carry; i++ {
		opList[i]++
		if opList[i] > opLast {
			opList[i] = day07OpAdd
		} else {
			carry = false
		}
	}
	return carry
}
