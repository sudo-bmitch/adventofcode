package main

import (
	"fmt"
	"io"
	"strings"
)

func init() {
	registerDay("09a", day09a)
	registerDay("09b", day09b)
}

func day09a(args []string, rdr io.Reader) (string, error) {
	inBytes, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	in := strings.TrimSpace(string(inBytes))
	posID := []int{}
	for i, r := range in {
		if r < '0' || r > '9' {
			return "", fmt.Errorf("unknown character in position %d: %c", i, r)
		}
		size := int(r - '0')
		add := make([]int, size)
		if i%2 == 0 {
			for j := range add {
				add[j] = i / 2
			}
		} else {
			for j := range add {
				add[j] = -1
			}
		}
		posID = append(posID, add...)
	}
	sum := 0
	// walk from each end, moving content from end to the beginning and adding the checksum
	tail := len(posID) - 1
	for i := range posID {
		if posID[i] == -1 {
			// move item from tail
			for posID[tail] < 0 && tail > i {
				tail--
			}
			if posID[tail] >= 0 && tail > i {
				posID[i], posID[tail] = posID[tail], -1
			}
		}
		if posID[i] < 0 {
			break // end of list
		}
		sum += (posID[i] * i)
	}
	return fmt.Sprintf("%d", sum), nil
}

func day09b(args []string, rdr io.Reader) (string, error) {
	inBytes, err := io.ReadAll(rdr)
	if err != nil {
		return "", err
	}
	in := strings.TrimSpace(string(inBytes))
	files := []day09File{}
	spaces := []day09File{}
	tail := 0
	for i, r := range in {
		if r < '0' || r > '9' {
			return "", fmt.Errorf("unknown character in position %d: %c", i, r)
		}
		size := int(r - '0')
		if i%2 == 0 {
			id := len(files)
			files = append(files, day09File{
				id:    id,
				start: tail,
				size:  size,
			})
		} else {
			spaces = append(spaces, day09File{
				start: tail,
				size:  size,
			})
		}
		tail += size
	}
	// walk files in reverse to see what can move
	for f := len(files) - 1; f > 0; f-- {
		s := 0
		for s < len(spaces) {
			if spaces[s].size >= files[f].size || spaces[s].start > files[f].start {
				break
			}
			s++
		}
		if s >= len(spaces) || spaces[s].start > files[f].start {
			continue // no space found, do not move this file
		}
		// move the file, shrink the space and start it after the file
		files[f].start = spaces[s].start
		spaces[s].size -= files[f].size
		spaces[s].start += files[f].size
	}
	// sum the result
	sum := 0
	for _, file := range files {
		sum += file.Sum()
	}
	return fmt.Sprintf("%d", sum), nil
}

type day09File struct {
	id    int
	start int
	size  int
}

func (f day09File) Sum() int {
	sum := 0
	for i := 0; i < f.size; i++ {
		sum += (i + f.start) * f.id
	}
	return sum
}
