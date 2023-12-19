package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type day19Op int

const (
	day19Gt day19Op = iota
	day19Lt
	day19None
)

type day19Step struct {
	param  rune
	op     day19Op
	value  int
	result string
}

type day19Workflow []day19Step

type day19Tool map[rune]int

type day19ToolSpans struct {
	minX, maxX, minM, maxM, minA, maxA, minS, maxS int
}

func day19a(args []string, rdr io.Reader) (string, error) {
	workflows, tools, err := day19Parse(rdr)
	if err != nil {
		return "", err
	}

	sum := 0
	for _, tool := range tools {
		curWf := "in"
		for curWf != "A" && curWf != "R" {
			if wf, ok := workflows[curWf]; ok {
				curWf, err = wf.GetResult(tool)
				if err != nil {
					return "", fmt.Errorf("failed running workflow %s for tool %v: %w", curWf, tool, err)
				}
			} else {
				return "", fmt.Errorf("unknown workflow %s", curWf)
			}
		}
		if curWf == "A" {
			// accepted
			sum += tool['x'] + tool['m'] + tool['a'] + tool['s']
		}
	}

	return fmt.Sprintf("%d", sum), nil
}

func day19b(args []string, rdr io.Reader) (string, error) {
	workflows, _, err := day19Parse(rdr)
	if err != nil {
		return "", err
	}
	ts := day19ToolSpans{
		minX: 1, maxX: 4000, minM: 1, maxM: 4000, minA: 1, maxA: 4000, minS: 1, maxS: 4000,
	}
	return fmt.Sprintf("%d", day19SumB(workflows, "in", ts)), nil
}

func day19Parse(rdr io.Reader) (map[string]day19Workflow, []day19Tool, error) {
	wfRe := regexp.MustCompile(`^([a-z]+)\{(.*)\}$`)
	toolRe := regexp.MustCompile(`^\{x=(\d+),m=(\d+),a=(\d+),s=(\d+)\}$`)
	opRe := regexp.MustCompile(`^([xmas])([><])(\d+):([a-zAR]+)$`)
	opLookup := map[string]day19Op{
		">": day19Gt,
		"<": day19Lt,
	}
	workflows := map[string]day19Workflow{}
	tools := []day19Tool{}
	in, err := io.ReadAll(rdr)
	if err != nil {
		return nil, nil, err
	}
	inSplit := strings.Split(strings.TrimSpace(string(in)), "\n\n")
	if len(inSplit) != 2 {
		return nil, nil, fmt.Errorf("input split is not 2 sections, found %d", len(inSplit))
	}
	for _, wfStr := range strings.Split(inSplit[0], "\n") {
		match := wfRe.FindStringSubmatch(wfStr)
		if len(match) != 3 {
			return nil, nil, fmt.Errorf("workflow expression could not be parsed (%d): %s", len(match), wfStr)
		}
		name := match[1]
		wf := day19Workflow{}
		for _, op := range strings.Split(match[2], ",") {
			opMatch := opRe.FindStringSubmatch(op)
			if len(opMatch) != 5 {
				wf = append(wf, day19Step{op: day19None, result: op})
				continue
			}
			opCmp, ok := opLookup[opMatch[2]]
			if !ok {
				return nil, nil, fmt.Errorf("failed to lookup op %s from workflow %s", opMatch[2], wfStr)
			}
			value, err := strconv.Atoi(opMatch[3])
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse value %s from workflow %s: %w", opMatch[3], wfStr, err)
			}
			wf = append(wf, day19Step{param: rune(opMatch[1][0]), op: opCmp, value: value, result: opMatch[4]})
		}
		workflows[name] = wf
	}

	for _, toolStr := range strings.Split(inSplit[1], "\n") {
		match := toolRe.FindStringSubmatch(toolStr)
		if len(match) != 5 {
			return nil, nil, fmt.Errorf("failed to parse tool (match len=%d): %s", len(match), toolStr)
		}
		x, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert x value (%s) to int on line %s: %w", match[1], toolStr, err)
		}
		m, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert m value (%s) to int on line %s: %w", match[2], toolStr, err)
		}
		a, err := strconv.Atoi(match[3])
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert a value (%s) to int on line %s: %w", match[3], toolStr, err)
		}
		s, err := strconv.Atoi(match[4])
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert s value (%s) to int on line %s: %w", match[4], toolStr, err)
		}
		tools = append(tools, day19Tool{'x': x, 'm': m, 'a': a, 's': s})
	}

	return workflows, tools, nil
}

func (wf day19Workflow) GetResult(tool day19Tool) (string, error) {
	for _, step := range wf {
		switch step.op {
		case day19Gt:
			if tool[step.param] > step.value {
				return step.result, nil
			}
		case day19Lt:
			if tool[step.param] < step.value {
				return step.result, nil
			}
		case day19None:
			return step.result, nil
		default:
			return "", fmt.Errorf("unhandled step op: %v", step.op)
		}
	}
	return "", fmt.Errorf("no default op in workflow")
}

func day19SumB(workflows map[string]day19Workflow, curWf string, toolSpans day19ToolSpans) int {
	if curWf == "R" {
		return 0
	}
	if curWf == "A" {
		if toolSpans.GtZero() {
			add := (toolSpans.maxX - toolSpans.minX + 1) *
				(toolSpans.maxM - toolSpans.minM + 1) *
				(toolSpans.maxA - toolSpans.minA + 1) *
				(toolSpans.maxS - toolSpans.minS + 1)
			// fmt.Fprintf(os.Stderr, "returning %d from accept with spans %v\n", add, toolSpans)
			return add
		}
	}
	wf, ok := workflows[curWf]
	if !ok {
		return 0 // not found, error?
	}
	sum := 0
	remainingSpan := toolSpans
	for _, step := range wf {
		switch step.op {
		case day19Gt:
			curSpan := remainingSpan
			// in min=1, max=2000, op > 1000, cur 1001-2000, remain 1-1000
			switch step.param {
			case 'x':
				if curSpan.minX < step.value+1 {
					curSpan.minX = step.value + 1
				}
				if remainingSpan.maxX > step.value {
					remainingSpan.maxX = step.value
				}
			case 'm':
				if curSpan.minM < step.value+1 {
					curSpan.minM = step.value + 1
				}
				if remainingSpan.maxM > step.value {
					remainingSpan.maxM = step.value
				}
			case 'a':
				if curSpan.minA < step.value+1 {
					curSpan.minA = step.value + 1
				}
				if remainingSpan.maxA > step.value {
					remainingSpan.maxA = step.value
				}
			case 's':
				if curSpan.minS < step.value+1 {
					curSpan.minS = step.value + 1
				}
				if remainingSpan.maxS > step.value {
					remainingSpan.maxS = step.value
				}
				// should default error?
			}
			if curSpan.GtZero() {
				add := day19SumB(workflows, step.result, curSpan)
				// fmt.Fprintf(os.Stderr, "from %s with spans %v adding %d\n", step.result, curSpan, add)
				sum += add
			}
		case day19Lt:
			curSpan := remainingSpan
			// in min=1, max=2000, op < 1000, cur 1-999, remain 1000-2000
			switch step.param {
			case 'x':
				if curSpan.maxX > step.value-1 {
					curSpan.maxX = step.value - 1
				}
				if remainingSpan.minX < step.value {
					remainingSpan.minX = step.value
				}
			case 'm':
				if curSpan.maxM > step.value-1 {
					curSpan.maxM = step.value - 1
				}
				if remainingSpan.minM < step.value {
					remainingSpan.minM = step.value
				}
			case 'a':
				if curSpan.maxA > step.value-1 {
					curSpan.maxA = step.value - 1
				}
				if remainingSpan.minA < step.value {
					remainingSpan.minA = step.value
				}
			case 's':
				if curSpan.maxS > step.value-1 {
					curSpan.maxS = step.value - 1
				}
				if remainingSpan.minS < step.value {
					remainingSpan.minS = step.value
				}
				// should default error?
			}
			if curSpan.GtZero() {
				add := day19SumB(workflows, step.result, curSpan)
				// fmt.Fprintf(os.Stderr, "from %s with spans %v adding %d\n", step.result, curSpan, add)
				sum += add
			}
		case day19None:
			if remainingSpan.GtZero() {
				add := day19SumB(workflows, step.result, remainingSpan)
				// fmt.Fprintf(os.Stderr, "from %s with spans %v adding %d\n", step.result, remainingSpan, add)
				sum += add
			}
			// should default error?
		}
	}
	return sum
}

func (ts day19ToolSpans) GtZero() bool {
	if ts.minX > ts.maxX || ts.minM > ts.maxM || ts.minA > ts.maxA || ts.minS > ts.maxS {
		return false
	}
	return true
}
