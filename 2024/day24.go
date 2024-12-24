package main

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func init() {
	registerDay("24a", day24a)
	registerDay("24b", day24b)
}

func day24a(args []string, rdr io.Reader) (string, error) {
	wires, err := day24ReadWires(rdr)
	if err != nil {
		return "", err
	}
	sum, err := day24SumVar('z', wires)
	return fmt.Sprintf("%d", sum), err
}

func day24b(args []string, rdr io.Reader) (string, error) {
	wires, err := day24ReadWires(rdr)
	if err != nil {
		return "", err
	}
	// expected addr gates:

	// half addr for 00:
	// x00 XOR y00 -> z00
	// x00 AND y00 -> cOut

	// full addr for ## > 0
	// x## XOR y## -> i##
	// cIn XOR i## -> z##
	// cIn AND i## -> j##
	// x## AND y## -> k##
	// j## OR k## -> cOut

	// loop through numbers, finding each required gate, and name of the output to search for next gate
	miswired := []string{}
	done := false
	carryIn, carryOut := "", ""
	for i := 0; !done; i++ {
		xName := fmt.Sprintf("x%02d", i)
		yName := fmt.Sprintf("y%02d", i)
		zName := fmt.Sprintf("z%02d", i)
		// if no inputs exist, this is the last entry
		if wires[xName] == nil || wires[yName] == nil {
			done = true
			// make sure carry == zName
			if carryIn != zName {
				miswired = append(miswired, carryIn)
			}
		} else if carryIn != "" {
			// full addr
			iName, jName, kName := "", "", ""
			zIn1, zIn2 := "", ""
			// first find the direct x/y/z connections
			for out, wire := range wires {
				// x## XOR y## -> i##
				if (wire.a == xName && wire.b == yName) || (wire.a == yName && wire.b == xName) {
					if wire.kind == day24Xor {
						iName = out
					}
					if wire.kind == day24And {
						kName = out
					}
				}
				// cIn XOR i## -> z##
				if out == zName && wire.kind == day24Xor {
					zIn1 = wire.a
					zIn2 = wire.b
				}
			}
			if iName == "" {
				fmt.Printf("could not find x%02d XOR y%02d gate\n", i, i)
			}
			if kName == "" {
				fmt.Printf("could not find x%02d AND y%02d gate\n", i, i)
			}
			// cIn XOR i## -> z##
			if zIn1 == "" || zIn2 == "" || (zIn1 != iName && zIn1 != carryIn && zIn2 != iName && zIn2 != carryIn) {
				miswired = append(miswired, zName)
			} else if zIn1 != iName && zIn2 != iName {
				miswired = append(miswired, iName)
				if zIn1 == carryIn {
					iName = zIn2
				} else {
					iName = zIn1
				}
				miswired = append(miswired, iName)
			} else if zIn1 != carryIn && zIn2 != carryIn {
				miswired = append(miswired, carryIn)
				if zIn1 == iName {
					carryIn = zIn2
				} else {
					carryIn = zIn1
				}
				miswired = append(miswired, carryIn)
			}
			// find the next hop
			for out, wire := range wires {
				// cIn XOR i## -> z## for something other than the expected zName
				if wire.kind == day24Xor && ((wire.a == carryIn && wire.b == iName) || (wire.a == iName && wire.b == carryIn)) && out != zName {
					miswired = append(miswired, out)
				}
				// c## AND i## -> j##
				if wire.kind == day24And && ((wire.a == carryIn && wire.b == iName) || (wire.a == iName && wire.b == carryIn)) {
					jName = out
				}
			}
			// j## OR k## -> cOut
			for out, wire := range wires {
				if wire.kind == day24Or && (wire.a == jName || wire.b == jName || wire.a == kName || wire.b == kName) {
					if wire.a != jName && wire.b != jName {
						miswired = append(miswired, jName)
					}
					if wire.a != kName && wire.b != kName {
						miswired = append(miswired, kName)
					}
					carryOut = out
				}
			}
		} else {
			// half addr
			foundZ := false
			foundC := false
			for out, wire := range wires {
				if (wire.a == xName && wire.b == yName) || (wire.a == yName && wire.b == xName) {
					if wire.kind == day24Xor {
						foundZ = true
						if out != zName {
							miswired = append(miswired, out)
						}
					}
					if wire.kind == day24And {
						foundC = true
						carryOut = out
					}
				}
				if foundZ && foundC {
					break
				}
			}
			if !foundZ || !foundC {
				fmt.Printf("could not find one of the gates for 00: foundZ = %t, foundC = %t\n", foundZ, foundC)
			}
		}
		carryIn = carryOut
	}

	sort.Strings(miswired)
	// dedup
	for j := len(miswired) - 1; j > 0; j-- {
		if miswired[j] == miswired[j-1] {
			if j == len(miswired)-1 {
				miswired = miswired[:len(miswired)-1]
			} else {
				miswired = append(miswired[:j], miswired[j+1:]...)
			}
		}
	}
	return strings.Join(miswired, ","), nil
}

func day24ReadWires(rdr io.Reader) (map[string]*day24Gate, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	inSplit := strings.Split(string(in), "\n\n")
	if len(inSplit) != 2 {
		return nil, fmt.Errorf("input does not have two parts, found %d", len(inSplit))
	}
	wires := map[string]*day24Gate{}
	for _, wireStr := range strings.Split(inSplit[0], "\n") {
		wireStr = strings.TrimSpace(wireStr)
		wireSplit := strings.Split(wireStr, ": ")
		if len(wireSplit) != 2 || (wireSplit[1] != "0" && wireSplit[1] != "1") {
			return nil, fmt.Errorf("unknown wire: %s", wireStr)
		}
		wires[wireSplit[0]] = &day24Gate{kind: day24Init, value: (wireSplit[1] == "1"), set: true}
	}
	gateRE := regexp.MustCompile(`^(\S+)\s+(AND|OR|XOR)\s+(\S+)\s+->\s+(\S+)$`)
	for _, gateStr := range strings.Split(inSplit[1], "\n") {
		gateStr = strings.TrimSpace(gateStr)
		if gateStr == "" {
			continue
		}
		gateMatch := gateRE.FindStringSubmatch(gateStr)
		if len(gateMatch) != 5 {
			return nil, fmt.Errorf("unknown gate (len=%d): %s", len(gateMatch), gateStr)
		}
		var kind day24GateKind
		switch gateMatch[2] {
		case "AND":
			kind = day24And
		case "OR":
			kind = day24Or
		case "XOR":
			kind = day24Xor
		default:
			return nil, fmt.Errorf("unknown gate (type=%s): %s", gateMatch[2], gateStr)
		}
		wires[gateMatch[4]] = &day24Gate{kind: kind, a: gateMatch[1], b: gateMatch[3]}
	}
	return wires, nil
}

func day24SumVar(prefix rune, wires map[string]*day24Gate) (int, error) {
	sum := 0
	for name := range wires {
		if name[0] == byte(prefix) {
			value, err := day24GetValue(name, wires)
			if err != nil {
				return 0, err
			}
			if value {
				num, err := strconv.Atoi(name[1:])
				if err != nil {
					return 0, err
				}
				inc := 1 << num
				sum += inc
			}
		}
	}
	return sum, nil
}

func day24GetValue(name string, wires map[string]*day24Gate) (bool, error) {
	return day24GetValueLoop(name, wires, []string{})
}

func day24GetValueLoop(name string, wires map[string]*day24Gate, parents []string) (bool, error) {
	if day24InList(name, parents) {
		return false, fmt.Errorf("loop detected")
	}
	parents = append(parents, name)
	if !wires[name].set {
		aValue, aErr := day24GetValueLoop(wires[name].a, wires, parents)
		if aErr != nil {
			return false, aErr
		}
		bValue, bErr := day24GetValueLoop(wires[name].b, wires, parents)
		if bErr != nil {
			return false, bErr
		}
		switch wires[name].kind {
		case day24And:
			wires[name].value = aValue && bValue
		case day24Or:
			wires[name].value = aValue || bValue
		case day24Xor:
			wires[name].value = (aValue || bValue) && !(aValue && bValue)
		default:
			panic(fmt.Errorf("unknown kind %d", wires[name].kind))
		}
		wires[name].set = true
	}
	return wires[name].value, nil
}

type day24Gate struct {
	a, b       string
	set, value bool
	kind       day24GateKind
}

type day24GateKind int

const (
	day24And day24GateKind = iota
	day24Or
	day24Xor
	day24Init
)

// type day24Addr struct {
// 	xyXor, xyAnd string
// 	zIn1, zIn2 string
// 	cIn, cOut string
// }

func day24InList(val string, list []string) bool {
	for _, cur := range list {
		if cur == val {
			return true
		}
	}
	return false
}
