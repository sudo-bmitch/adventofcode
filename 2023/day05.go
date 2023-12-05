package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func day05a(args []string, rdr io.Reader) (string, error) {
	lowest := -1
	a, err := day05Parse(rdr)
	if err != nil {
		return "", err
	}
	for _, seed := range a.seeds {
		loc := a.seedToLoc(seed)
		if lowest < 0 || loc < lowest {
			lowest = loc
		}
	}
	return fmt.Sprintf("%d", lowest), nil
}

func day05b(args []string, rdr io.Reader) (string, error) {
	lowest := -1
	a, err := day05Parse(rdr)
	if err != nil {
		return "", err
	}
	for seedRange := 0; seedRange+1 < len(a.seeds); seedRange += 2 {
		seedStart := a.seeds[seedRange]
		length := a.seeds[seedRange+1]
		for i := 0; i < length; i++ {
			loc, skip := a.seedToLocWithMaxSkip(seedStart + i)
			if lowest < 0 || loc < lowest {
				lowest = loc
			}
			if skip > 0 {
				i += skip
			}
		}
	}
	return fmt.Sprintf("%d", lowest), nil
}

func day05Parse(rdr io.Reader) (day05Almanac, error) {
	a := day05Almanac{}
	in := bufio.NewScanner(rdr)
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[:7] != "seeds: " {
			return a, fmt.Errorf("seeds not found on %s", line)
		}
		seeds, err := parseNumListBySpace(line[7:])
		if err != nil {
			return a, fmt.Errorf("failed to parse seed list from %s: %w", line, err)
		}
		a.seeds = seeds
		break
	}
	if !in.Scan() || in.Text() != "" {
		return a, fmt.Errorf("missing whitespace after seeds")
	}
	if !in.Scan() || in.Text() != "seed-to-soil map:" {
		return a, fmt.Errorf("missing seed-to-soil map")
	}
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		convert, err := parseNumListBySpace(line)
		if err != nil || len(convert) != 3 {
			return a, fmt.Errorf("failed to parse seed-to-soil entry: %s", line)
		}
		a.seedToSoil.entries = append(a.seedToSoil.entries, day05Convert{dest: convert[0], src: convert[1], length: convert[2]})
	}
	if !in.Scan() || in.Text() != "soil-to-fertilizer map:" {
		return a, fmt.Errorf("missing soil-to-fertilizer map")
	}
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		convert, err := parseNumListBySpace(line)
		if err != nil || len(convert) != 3 {
			return a, fmt.Errorf("failed to parse soil-to-fertilizer entry: %s", line)
		}
		a.soilToFert.entries = append(a.soilToFert.entries, day05Convert{dest: convert[0], src: convert[1], length: convert[2]})
	}
	if !in.Scan() || in.Text() != "fertilizer-to-water map:" {
		return a, fmt.Errorf("missing fertilizer-to-water map")
	}
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		convert, err := parseNumListBySpace(line)
		if err != nil || len(convert) != 3 {
			return a, fmt.Errorf("failed to parse fertilizer-to-water entry: %s", line)
		}
		a.fertToWater.entries = append(a.fertToWater.entries, day05Convert{dest: convert[0], src: convert[1], length: convert[2]})
	}
	if !in.Scan() || in.Text() != "water-to-light map:" {
		return a, fmt.Errorf("missing water-to-light map")
	}
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		convert, err := parseNumListBySpace(line)
		if err != nil || len(convert) != 3 {
			return a, fmt.Errorf("failed to parse water-to-light entry: %s", line)
		}
		a.waterToLight.entries = append(a.waterToLight.entries, day05Convert{dest: convert[0], src: convert[1], length: convert[2]})
	}
	if !in.Scan() || in.Text() != "light-to-temperature map:" {
		return a, fmt.Errorf("missing light-to-temperature map")
	}
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		convert, err := parseNumListBySpace(line)
		if err != nil || len(convert) != 3 {
			return a, fmt.Errorf("failed to parse light-to-temperature entry: %s", line)
		}
		a.lightToTemp.entries = append(a.lightToTemp.entries, day05Convert{dest: convert[0], src: convert[1], length: convert[2]})
	}
	if !in.Scan() || in.Text() != "temperature-to-humidity map:" {
		return a, fmt.Errorf("missing temperature-to-humidity map")
	}
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		convert, err := parseNumListBySpace(line)
		if err != nil || len(convert) != 3 {
			return a, fmt.Errorf("failed to parse temperature-to-humidity entry: %s", line)
		}
		a.tempToHumid.entries = append(a.tempToHumid.entries, day05Convert{dest: convert[0], src: convert[1], length: convert[2]})
	}
	if !in.Scan() || in.Text() != "humidity-to-location map:" {
		return a, fmt.Errorf("missing humidity-to-location map")
	}
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		convert, err := parseNumListBySpace(line)
		if err != nil || len(convert) != 3 {
			return a, fmt.Errorf("failed to parse humidity-to-location entry: %s", line)
		}
		a.humidToLoc.entries = append(a.humidToLoc.entries, day05Convert{dest: convert[0], src: convert[1], length: convert[2]})
	}

	return a, nil
}

type day05Almanac struct {
	seeds        []int
	seedToSoil   day05ConvertSet
	soilToFert   day05ConvertSet
	fertToWater  day05ConvertSet
	waterToLight day05ConvertSet
	lightToTemp  day05ConvertSet
	tempToHumid  day05ConvertSet
	humidToLoc   day05ConvertSet
}

func (a day05Almanac) seedToLoc(seed int) int {
	soil := a.seedToSoil.dest(seed)
	fert := a.soilToFert.dest(soil)
	water := a.fertToWater.dest(fert)
	light := a.waterToLight.dest(water)
	temp := a.lightToTemp.dest(light)
	humid := a.tempToHumid.dest(temp)
	return a.humidToLoc.dest(humid)
}

func (a day05Almanac) seedToLocWithMaxSkip(seed int) (int, int) {
	soil, skip := a.seedToSoil.destWithSkip(seed, -1)
	fert, skip := a.soilToFert.destWithSkip(soil, skip)
	water, skip := a.fertToWater.destWithSkip(fert, skip)
	light, skip := a.waterToLight.destWithSkip(water, skip)
	temp, skip := a.lightToTemp.destWithSkip(light, skip)
	humid, skip := a.tempToHumid.destWithSkip(temp, skip)
	return a.humidToLoc.destWithSkip(humid, skip)
}

type day05ConvertSet struct {
	entries []day05Convert
}

type day05Convert struct {
	src, dest, length int
}

func (s day05ConvertSet) dest(src int) int {
	for _, e := range s.entries {
		if src >= e.src && src < e.src+e.length {
			return e.dest + (src - e.src)
		}
	}
	return src
}

func (s day05ConvertSet) destWithSkip(src, skip int) (int, int) {
	for _, e := range s.entries {
		if src >= e.src && src < e.src+e.length {
			if skip < 0 || skip > e.length-(src-e.src)-1 {
				skip = e.length - (src - e.src) - 1
			}
			return e.dest + (src - e.src), skip
		}
		if src < e.src && (skip < 0 || skip > e.src-src-1) {
			skip = e.src - src - 1
		}
	}
	return src, skip
}
