package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func init() {
	registerDay("17a", day17a)
	registerDay("17b", day17b)
	registerDay("17slow", day17slow)
}

func day17a(args []string, rdr io.Reader) (string, error) {
	vm := day17VM{}
	err := vm.Load(rdr)
	if err != nil {
		return "", err
	}
	err = vm.Run()
	if err != nil {
		return "", err
	}
	return vm.PrintOut(), nil
}

func day17slow(args []string, rdr io.Reader) (string, error) {
	vm := day17VM{}
	err := vm.Load(rdr)
	if err != nil {
		return "", err
	}
	vm.outRequire = vm.program
	regA := 0
	regB := vm.regB
	regC := vm.regC
	for {
		// reset program
		vm.pc, vm.regA, vm.regB, vm.regC = 0, regA, regB, regC
		vm.out = []int{}
		// run to see if desired output is returned
		err = vm.Run()
		if err == nil && len(vm.out) == len(vm.outRequire) {
			break
		}
		regA++
		if regA%1000000 == 0 {
			fmt.Printf("tested %d programs\n", regA)
		}
	}
	return fmt.Sprintf("%d", regA), nil
}

func day17b(args []string, rdr io.Reader) (string, error) {
	vm := day17VM{}
	err := vm.Load(rdr)
	if err != nil {
		return "", err
	}
	regA := 0
	regB := vm.regB
	regC := vm.regC
	// each chunk is a 3 bit int (0-7), incremented/appended from the tail to generate the tail of the program
	regAChunk := []int{0}
	for {
		// assemble regA
		regA = 0
		for _, chunk := range regAChunk {
			regA <<= 3
			regA += chunk
		}
		// reset program
		vm.pc, vm.regA, vm.regB, vm.regC = 0, regA, regB, regC
		vm.out = []int{}
		// compute the desired tail of the program
		vm.outRequire = vm.program[len(vm.program)-len(regAChunk):]
		// run
		err = vm.Run()
		if err == nil && len(vm.out) == len(vm.outRequire) {
			if len(vm.out) == len(vm.program) {
				// done
				break
			} else {
				fmt.Printf("success regA: Oo%o, output: %v\n", regA, vm.out)
				// chunk works, append another
				regAChunk = append(regAChunk, 0)
			}
		} else {
			fmt.Printf("failed regA: Oo%o, output: %v\n", regA, vm.out)
			// chunk does not work, increment
			for len(regAChunk) > 0 {
				regAChunk[len(regAChunk)-1]++
				if regAChunk[len(regAChunk)-1] < 8 {
					break
				}
				// pop
				regAChunk = regAChunk[:len(regAChunk)-1]
			}
		}
	}
	fmt.Printf("program %v outputs %v with regA %d[Oo%o]\n", vm.program, vm.out, regA, regA)
	return fmt.Sprintf("%d", regA), nil
}

type day17VM struct {
	out              []int
	pc               int
	regA, regB, regC int
	program          []int
	outRequire       []int
}

type day17HaltCheck struct {
	pc, regA, regB, regC int
}

func (vm *day17VM) Load(rdr io.Reader) error {
	input, err := io.ReadAll(rdr)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(input), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		split := strings.SplitN(line, ": ", 2)
		if len(split) == 1 {
			return fmt.Errorf("unknown line (missing colon): %s", line)
		}
		switch split[0] {
		case "Register A":
			i, err := strconv.Atoi(split[1])
			if err != nil {
				return err
			}
			vm.regA = i
		case "Register B":
			i, err := strconv.Atoi(split[1])
			if err != nil {
				return err
			}
			vm.regB = i
		case "Register C":
			i, err := strconv.Atoi(split[1])
			if err != nil {
				return err
			}
			vm.regC = i
		case "Program":
			splitP := strings.Split(split[1], ",")
			for _, progStr := range splitP {
				i, err := strconv.Atoi(progStr)
				if err != nil {
					return err
				}
				vm.program = append(vm.program, i)
			}
		default:
			return fmt.Errorf("unknown line (first param): %s", line)
		}
	}
	return nil
}

func (vm *day17VM) PrintOut() string {
	outS := make([]string, len(vm.out))
	for i, v := range vm.out {
		outS[i] = strconv.Itoa(v)
	}
	return strings.Join(outS, ",")
}

func (vm *day17VM) Run() error {
	haltCheck := map[day17HaltCheck]bool{}
	for vm.pc < len(vm.program) {
		hc := day17HaltCheck{pc: vm.pc, regA: vm.regA, regB: vm.regB, regC: vm.regC}
		if haltCheck[hc] {
			return fmt.Errorf("program does not halt")
		}
		haltCheck[hc] = true
		pcNext := vm.pc + 2
		switch vm.program[vm.pc] {
		case 0: // adv
			val, err := vm.comboVal()
			if err != nil {
				return err
			}
			vm.regA >>= val
		case 1: // bxl
			if len(vm.program) <= vm.pc+1 {
				return fmt.Errorf("read bxl value past end of program, pc=%d", vm.pc)
			}
			vm.regB ^= vm.program[vm.pc+1]
		case 2: // bst
			val, err := vm.comboVal()
			if err != nil {
				return err
			}
			vm.regB = val & 0x7
		case 3: // jnz
			if vm.regA != 0 {
				if len(vm.program) <= vm.pc+1 {
					return fmt.Errorf("read jump value past end of program, pc=%d", vm.pc)
				}
				pcNext = vm.program[vm.pc+1]
			}
		case 4: // bxc
			vm.regB ^= vm.regC
		case 5: // out
			val, err := vm.comboVal()
			if err != nil {
				return err
			}
			vm.out = append(vm.out, val&0x7)
			if vm.outRequire != nil {
				if len(vm.out) > len(vm.outRequire) {
					return fmt.Errorf("output exceeded max length, %v > %v", vm.out, vm.outRequire)
				}
				if vm.out[len(vm.out)-1] != vm.outRequire[len(vm.out)-1] {
					return fmt.Errorf("output mismatch, %v != %v", vm.out, vm.outRequire)
				}
			}
		case 6: // bdv
			val, err := vm.comboVal()
			if err != nil {
				return err
			}
			vm.regB = vm.regA >> val
		case 7: // cdv
			val, err := vm.comboVal()
			if err != nil {
				return err
			}
			vm.regC = vm.regA >> val
		default:
			return fmt.Errorf("unhanded instruction %d at pc %d", vm.program[vm.pc], vm.pc)
		}
		vm.pc = pcNext
	}
	return nil
}

func (vm *day17VM) comboVal() (int, error) {
	if len(vm.program) <= vm.pc+1 {
		return 0, fmt.Errorf("read combo value past end of program, pc=%d", vm.pc)
	}
	combo := vm.program[vm.pc+1]
	switch combo {
	case 0, 1, 2, 3:
		return combo, nil
	case 4:
		return vm.regA, nil
	case 5:
		return vm.regB, nil
	case 6:
		return vm.regC, nil
	}
	return 0, fmt.Errorf("unsupported combo value %d", combo)
}
