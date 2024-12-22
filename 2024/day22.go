package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func init() {
	registerDay("22a", day22a)
	registerDay("22b", day22b)
}

func day22a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(in), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		num, err := strconv.Atoi(line)
		if err != nil {
			return "", fmt.Errorf("failed to convert line %s: %w", line, err)
		}
		rnd := day22Rnd(num)
		for range 2000 {
			rnd.next()
		}
		fmt.Printf("Seed %d generates %d\n", num, int(rnd))
		sum += int(rnd)
	}

	return fmt.Sprintf("%d", sum), nil
}

func day22b(args []string, rdr io.Reader) (string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	sellerSeq := map[string][]*int{}
	lines := strings.Split(string(in), "\n")
	sellerCount := len(lines)
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		num, err := strconv.Atoi(line)
		if err != nil {
			return "", fmt.Errorf("failed to convert line %s: %w", line, err)
		}
		rnd := day22Rnd(num)
		priceList := []int{num % 10}
		for range 2000 {
			rnd.next()
			priceList = append(priceList, int(rnd)%10)
			if len(priceList) == 5 {
				seq := fmt.Sprintf("%d,%d,%d,%d", priceList[1]-priceList[0], priceList[2]-priceList[1], priceList[3]-priceList[2], priceList[4]-priceList[3])
				priceList = priceList[1:]
				if _, ok := sellerSeq[seq]; !ok {
					sellerSeq[seq] = make([]*int, sellerCount)
				}
				if sellerSeq[seq][i] == nil {
					price := priceList[3]
					sellerSeq[seq][i] = &price
				}
			}
		}
	}

	bestSeq := ""
	bestSale := 0
	for seq := range sellerSeq {
		curSale := 0
		for _, price := range sellerSeq[seq] {
			if price != nil {
				curSale += *price
			}
		}
		if curSale > bestSale {
			bestSeq = seq
			bestSale = curSale
		}
	}
	fmt.Printf("For sequence %s, total sale is %d:\n", bestSeq, bestSale)
	return fmt.Sprintf("%d", bestSale), nil
}

type day22Rnd int

func (r *day22Rnd) next() int {
	*r = ((*r << 6) ^ *r) % 16777216
	*r = ((*r >> 5) ^ *r) % 16777216
	*r = ((*r << 11) ^ *r) % 16777216
	return int(*r)
}
