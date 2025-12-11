package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func init() {
	registerDay("11a", day11a)
	registerDay("11b", day11b)
}

func day11a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	debug := false
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	type day11aInputs struct {
		count int
		devs  []string
	}
	devices := map[string][]string{}
	inputs := map[string]*day11aInputs{
		"you": {
			count: 1,
		},
	}
	for line := range strings.SplitSeq(strings.TrimSpace(string(in)), "\n") {
		fields := strings.Split(line, " ")
		if fields[0][len(fields[0])-1] != ':' {
			return "", fmt.Errorf("line missing colon after first field: %s", line)
		}
		dev := fields[0][:len(fields[0])-1]
		outs := fields[1:]
		devices[dev] = outs
		for _, out := range outs {
			if inputs[out] == nil {
				inputs[out] = &day11aInputs{count: -1}
			}
			inputs[out].devs = append(inputs[out].devs, dev)
		}
	}
	for inputs["out"].count < 0 {
		progress := false
		for inDev, input := range inputs {
			if input.count >= 0 {
				continue
			}
			complete := true
			count := 0
			for _, outDev := range input.devs {
				if inputs[outDev] == nil {
					continue
				} else if inputs[outDev].count >= 0 {
					count += inputs[outDev].count
				} else {
					complete = false
					break
				}
			}
			if !complete {
				continue
			}
			progress = true
			input.count = count
			if debug {
				fmt.Fprintf(os.Stderr, "device %s has %d input paths\n", inDev, count)
			}
		}
		if !progress {
			return "", fmt.Errorf("failed to make any progress, loop on an input?")
		}
		if debug {
			fmt.Fprintf(os.Stderr, "end of a pass\n")
		}
	}
	sum = inputs["out"].count
	return fmt.Sprintf("%d", sum), nil
}

func day11b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	debug := false
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	type day11aInputs struct {
		countNeither, countDAC, countFFT, countBoth int
		devs                                        []string
	}
	devices := map[string][]string{}
	inputs := map[string]*day11aInputs{
		"svr": {
			countNeither: 1,
		},
	}
	for line := range strings.SplitSeq(strings.TrimSpace(string(in)), "\n") {
		fields := strings.Split(line, " ")
		if fields[0][len(fields[0])-1] != ':' {
			return "", fmt.Errorf("line missing colon after first field: %s", line)
		}
		dev := fields[0][:len(fields[0])-1]
		outs := fields[1:]
		devices[dev] = outs
		for _, out := range outs {
			if inputs[out] == nil {
				inputs[out] = &day11aInputs{countNeither: -1, countDAC: -1, countFFT: -1, countBoth: -1}
			}
			inputs[out].devs = append(inputs[out].devs, dev)
		}
	}
	for inputs["out"].countBoth < 0 {
		progress := false
		for inDev, input := range inputs {
			if input.countNeither >= 0 {
				continue
			}
			complete := true
			countNeither, countDAC, countFFT, countBoth := 0, 0, 0, 0
			for _, outDev := range input.devs {
				if inputs[outDev] == nil {
					continue
				} else if inputs[outDev].countNeither >= 0 {
					countNeither += inputs[outDev].countNeither
					countDAC += inputs[outDev].countDAC
					countFFT += inputs[outDev].countFFT
					countBoth += inputs[outDev].countBoth
				} else {
					complete = false
					break
				}
			}
			if !complete {
				continue
			}
			progress = true
			if inDev == "dac" {
				countBoth += countFFT
				countFFT = 0
				countDAC += countNeither
				countNeither = 0
			}
			if inDev == "fft" {
				countBoth += countDAC
				countDAC = 0
				countFFT += countNeither
				countNeither = 0
			}
			input.countNeither = countNeither
			input.countDAC = countDAC
			input.countFFT = countFFT
			input.countBoth = countBoth
			if debug {
				fmt.Fprintf(os.Stderr, "device %s has %d, %d, %d, %d input paths\n", inDev, countNeither, countDAC, countFFT, countBoth)
			}
		}
		if !progress {
			return "", fmt.Errorf("failed to make any progress, loop on an input?")
		}
		if debug {
			fmt.Fprintf(os.Stderr, "end of a pass\n")
		}
	}
	sum = inputs["out"].countBoth
	return fmt.Sprintf("%d", sum), nil
}
