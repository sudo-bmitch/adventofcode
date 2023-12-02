package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type day02Game struct {
	number int
	hands  []day02Group
}

type day02Group struct {
	red, green, blue int
}

func day02a(args []string, rdr io.Reader) (string, error) {
	sum := 0
	possibleBag := day02Group{
		red:   12,
		green: 13,
		blue:  14,
	}
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		g, err := day02Parse(line)
		if err != nil {
			return "", err
		}
		impossible := false
		for _, hand := range g.hands {
			if hand.red > possibleBag.red || hand.green > possibleBag.green || hand.blue > possibleBag.blue {
				impossible = true
				break
			}
		}
		if impossible {
			continue
		}
		sum += g.number
	}
	return fmt.Sprintf("%d", sum), nil
}

func day02b(args []string, rdr io.Reader) (string, error) {
	sum := 0
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		g, err := day02Parse(line)
		if err != nil {
			return "", err
		}
		minReq := day02Group{}
		for _, hand := range g.hands {
			if hand.red > minReq.red {
				minReq.red = hand.red
			}
			if hand.blue > minReq.blue {
				minReq.blue = hand.blue
			}
			if hand.green > minReq.green {
				minReq.green = hand.green
			}
		}
		powerSet := minReq.red * minReq.blue * minReq.green
		sum += powerSet
	}
	return fmt.Sprintf("%d", sum), nil
}

func day02Parse(line string) (day02Game, error) {
	g := day02Game{
		hands: []day02Group{},
	}
	gameSplit := strings.SplitN(line, ": ", 2)
	if len(gameSplit) < 2 || gameSplit[0][:5] != "Game " {
		return g, fmt.Errorf("colon separator or game prefix missing in %s", line)
	}
	n, err := strconv.Atoi(gameSplit[0][5:])
	if err != nil {
		return g, fmt.Errorf("unable to parse game number %s from %s", gameSplit[0][5:], line)
	}
	g.number = n
	for _, handStr := range strings.Split(gameSplit[1], "; ") {
		curHand := day02Group{}
		for _, colorStr := range strings.Split(handStr, ", ") {
			colorSplit := strings.Split(colorStr, " ")
			if len(colorSplit) != 2 {
				return g, fmt.Errorf("color split does not have two fields: %s, from line: %s", colorStr, line)
			}
			n, err := strconv.Atoi(colorSplit[0])
			if err != nil {
				return g, fmt.Errorf("failed to parse number from %s, line %s: %w", colorSplit[0], line, err)
			}
			switch colorSplit[1] {
			case "red":
				curHand.red = n
			case "green":
				curHand.green = n
			case "blue":
				curHand.blue = n
			default:
				return g, fmt.Errorf("unknown color %s from hand %s in line %s", colorSplit[1], handStr, line)
			}
		}
		g.hands = append(g.hands, curHand)
	}

	return g, nil
}
