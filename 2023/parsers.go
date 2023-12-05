package main

import (
	"strconv"
	"strings"
)

func parseNumListBySpace(numStr string) ([]int, error) {
	numList := []int{}
	numSplit := strings.Split(numStr, " ")
	for _, curStr := range numSplit {
		if curStr != "" {
			i, err := strconv.Atoi(curStr)
			if err != nil {
				return numList, err
			}
			numList = append(numList, i)
		}
	}
	return numList, nil
}
