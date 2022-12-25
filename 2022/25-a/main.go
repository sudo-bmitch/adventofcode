package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

var debug = true

type base5 int
type snafu []base5

func main() {
	in := bufio.NewScanner(os.Stdin)
	lineRE := regexp.MustCompile(`^[=\-012]+$`)
	sum := 0
	// parse input
	for in.Scan() {
		line := in.Text()
		if line == "" {
			continue
		}
		if !lineRE.MatchString(line) {
			fmt.Fprintf(os.Stderr, "failed to parse %s\n", line)
			return
		}
		sVal := MustStrToSnafu(line)
		dVal := sVal.Int()
		sum += dVal
		if debug {
			fmt.Printf("debug: string %s, snafu %s, dec %d, sum %d\n", line, sVal.String(), dVal, sum)
		}
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v\n", err)
		return
	}

	result := MustDecToSnafu(sum)
	fmt.Printf("Result: %d = %s\n", sum, result.String())
}

func (b base5) Byte() byte {
	switch b {
	case -2:
		return '='
	case -1:
		return '-'
	case 0:
		return '0'
	case 1:
		return '1'
	case 2:
		return '2'
	default:
		return '?'
	}
}
func MustByteToBase5(c byte) base5 {
	switch c {
	case '=':
		return -2
	case '-':
		return -1
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	default:
		panic(fmt.Sprintf("invalid base5 value %c", c))
	}
}

func MustStrToSnafu(s string) snafu {
	sl := []base5{}
	for i := len(s) - 1; i >= 0; i-- {
		sl = append(sl, MustByteToBase5(s[i]))
	}
	return sl
}

func MustDecToSnafu(dec int) snafu {
	sl := []base5{}
	cur := dec
	for cur != 0 {
		rem := ((cur + 2) % 5) - 2
		cur = (cur - rem) / 5
		sl = append(sl, base5(rem))
	}
	return sl
}

func (s snafu) String() string {
	bl := []byte{}
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		bl = append(bl, c.Byte())
	}
	return string(bl)
}

func (s snafu) Int() int {
	val := 0
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		val = val*5 + int(c)
	}
	return val
}

// func min(a, b int) int {
// 	if a > b {
// 		return b
// 	}
// 	return a
// }
// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// func mustAtoi(s string) int {
// 	i, err := strconv.Atoi(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return i
// }
