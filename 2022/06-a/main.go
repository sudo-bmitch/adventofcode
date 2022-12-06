package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const uniqueLen = 4

func main() {
	inB, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read input: %v", err)
		return
	}
	in := string(inB)
	in = strings.TrimSpace(in)
	uniqBuf := [uniqueLen]rune{}
	for i, c := range in {
		if i < uniqueLen-1 {
			uniqBuf[i] = c
			continue
		}
		uniqBuf[i%uniqueLen] = c
		if isUnique(uniqBuf) {
			fmt.Printf("found unique range of %d ending at %d\n", uniqueLen, i+1)
			break
		}
	}
}

func isUnique(r [uniqueLen]rune) bool {
	for start := 0; start < uniqueLen-1; start++ {
		for offset := start + 1; offset < uniqueLen; offset++ {
			if r[start] == r[offset] {
				return false
			}
		}
	}
	return true
}
