package parse

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"strconv"
	"strings"
)

// MustNumSlice iterates over an input of white-space separated numbers.
// Failure to parse a string into a number will panic.
func MustNumSlice(rdr io.Reader) iter.Seq[[]int] {
	in := bufio.NewScanner(rdr)
	return func(yield func([]int) bool) {
		var err error
		for in.Scan() {
			line := in.Text()
			splitStr := strings.Fields(line)
			if len(splitStr) == 0 {
				continue
			}
			splitInt := make([]int, len(splitStr))
			for i, s := range splitStr {
				splitInt[i], err = strconv.Atoi(s)
				if err != nil {
					panic(fmt.Sprintf("could not parse number \"%s\": %v", s, err))
				}
			}
			if !yield(splitInt) {
				return
			}
		}
	}
}
