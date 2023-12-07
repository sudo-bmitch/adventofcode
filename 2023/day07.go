package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

type day07HandType int

const (
	day07HTHighest = iota
	day07HTOnePair
	day07HTTwoPair
	day07HTThree
	day07HTFullHouse
	day07HTFour
	day07HTFive
)

var day07aCardVal = map[rune]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

var day07bCardVal = map[rune]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 1,
	'Q': 12,
	'K': 13,
	'A': 14,
}

type day07Bid struct {
	cards  []rune
	amount int
	ht     day07HandType
}

func day07a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	bids := []day07Bid{}
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		b, err := day07Parse(line)
		if err != nil {
			return "", err
		}
		err = b.day07aEvaluate()
		if err != nil {
			return "", err
		}
		bids = append(bids, b)
	}
	slices.SortFunc(bids, day07aBidCmp)
	for i, b := range bids {
		// fmt.Fprintf(os.Stderr, "rank %d, hand %s, bid %d, winning %d\n", i+1, string(b.cards), b.amount, b.amount*(i+1))
		sum += b.amount * (i + 1)
	}
	return fmt.Sprintf("%d", sum), nil
}

func day07b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	bids := []day07Bid{}
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		b, err := day07Parse(line)
		if err != nil {
			return "", err
		}
		err = b.day07bEvaluate()
		if err != nil {
			return "", err
		}
		bids = append(bids, b)
	}
	slices.SortFunc(bids, day07bBidCmp)
	for i, b := range bids {
		// fmt.Fprintf(os.Stderr, "rank %d, hand %s, bid %d, winning %d\n", i+1, string(b.cards), b.amount, b.amount*(i+1))
		sum += b.amount * (i + 1)
	}
	return fmt.Sprintf("%d", sum), nil
}

func day07Parse(line string) (day07Bid, error) {
	b := day07Bid{}
	lSplit := strings.Split(line, " ")
	if len(lSplit) != 2 || len(lSplit[0]) != 5 {
		return b, fmt.Errorf("invalid line: %s", line)
	}
	b.cards = []rune(lSplit[0])
	i, err := strconv.Atoi(lSplit[1])
	if err != nil {
		return b, fmt.Errorf("failed parsing amount in %s: %w", line, err)
	}
	b.amount = i
	return b, nil
}

func (b *day07Bid) day07aEvaluate() error {
	cards := map[rune]int{}
	for _, r := range b.cards {
		cards[r]++
	}
	counts := []int{}
	for _, i := range cards {
		counts = append(counts, i)
	}
	slices.Sort(counts)
	slices.Reverse(counts)
	switch {
	case len(counts) == 1:
		b.ht = day07HTFive
	case len(counts) == 2 && counts[0] == 4:
		b.ht = day07HTFour
	case len(counts) == 2 && counts[0] == 3:
		b.ht = day07HTFullHouse
	case len(counts) == 3 && counts[0] == 3:
		b.ht = day07HTThree
	case len(counts) == 3 && counts[0] == 2:
		b.ht = day07HTTwoPair
	case len(counts) == 4 && counts[0] == 2:
		b.ht = day07HTOnePair
	case len(counts) == 5:
		b.ht = day07HTHighest
	default:
		return fmt.Errorf("unhandled hand: %s", string(b.cards))
	}
	return nil
}

func day07aBidCmp(a, b day07Bid) int {
	if a.ht < b.ht {
		return -1
	}
	if a.ht > b.ht {
		return 1
	}
	for i := range a.cards {
		if day07aCardVal[a.cards[i]] < day07aCardVal[b.cards[i]] {
			return -1
		}
		if day07aCardVal[a.cards[i]] > day07aCardVal[b.cards[i]] {
			return 1
		}
	}
	return 0
}

func (b *day07Bid) day07bEvaluate() error {
	cards := map[rune]int{}
	jokers := 0
	for _, r := range b.cards {
		if r == 'J' {
			jokers++
		} else {
			cards[r]++
		}
	}
	counts := []int{}
	for _, i := range cards {
		counts = append(counts, i)
	}
	slices.Sort(counts)
	slices.Reverse(counts)
	if len(counts) > 0 {
		counts[0] += jokers
	} else {
		counts = []int{5} // all jokers
	}
	switch {
	case len(counts) == 1:
		b.ht = day07HTFive
	case len(counts) == 2 && counts[0] == 4:
		b.ht = day07HTFour
	case len(counts) == 2 && counts[0] == 3:
		b.ht = day07HTFullHouse
	case len(counts) == 3 && counts[0] == 3:
		b.ht = day07HTThree
	case len(counts) == 3 && counts[0] == 2:
		b.ht = day07HTTwoPair
	case len(counts) == 4 && counts[0] == 2:
		b.ht = day07HTOnePair
	case len(counts) == 5:
		b.ht = day07HTHighest
	default:
		return fmt.Errorf("unhandled hand: %s", string(b.cards))
	}
	return nil
}

func day07bBidCmp(a, b day07Bid) int {
	if a.ht < b.ht {
		return -1
	}
	if a.ht > b.ht {
		return 1
	}
	for i := range a.cards {
		if day07bCardVal[a.cards[i]] < day07bCardVal[b.cards[i]] {
			return -1
		}
		if day07bCardVal[a.cards[i]] > day07bCardVal[b.cards[i]] {
			return 1
		}
	}
	return 0
}
