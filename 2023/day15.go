package main

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type day15Hash struct {
	h int
}

type day15Box struct {
	l []day15Lens
}

type day15Lens struct {
	label string
	focal int
}

func day15a(args []string, rdr io.Reader) (string, error) {
	seqList, err := day15Parse(rdr)
	if err != nil {
		return "", err
	}
	sum := 0
	for _, seq := range seqList {
		h := day15Hash{}
		_, err := h.Write([]byte(seq))
		if err != nil {
			return "", err
		}
		sum += h.Sum()
	}

	return fmt.Sprintf("%d", sum), nil
}

func day15b(args []string, rdr io.Reader) (string, error) {
	boxes := make([]day15Box, 256)
	seqList, err := day15Parse(rdr)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`^([a-z]+)(\-|=([0-9]))$`)
	for _, seq := range seqList {
		if len(seq) < 3 {
			return "", fmt.Errorf("invalid sequence %s", seq)
		}
		match := re.FindStringSubmatch(seq)
		if len(match) < 3 {
			return "", fmt.Errorf("invalid sequence %s, regexp match %v", seq, match)
		}
		label := match[1]
		op := match[2][0]
		h := day15Hash{}
		_, err = h.Write([]byte(label))
		if err != nil {
			return "", err
		}
		box := h.Sum()
		if boxes[box].l == nil {
			boxes[box].l = []day15Lens{}
		}
		switch op {
		case '-':
			// search for label and remove
			for i, lens := range boxes[box].l {
				if lens.label == label {
					boxes[box].l = append(boxes[box].l[:i], boxes[box].l[i+1:]...)
					break
				}
			}
		case '=':
			if len(match) != 4 {
				return "", fmt.Errorf("invalid sequence, missing focal length in %s", seq)
			}
			focal, err := strconv.Atoi(match[3])
			if err != nil {
				return "", fmt.Errorf("invalid focal in %s: %w", seq, err)
			}
			found := false
			for i, lens := range boxes[box].l {
				if lens.label == label {
					boxes[box].l[i].focal = focal
					found = true
					break
				}
			}
			if !found {
				boxes[box].l = append(boxes[box].l, day15Lens{label: label, focal: focal})
			}
		default:
			return "", fmt.Errorf("unknown operation %c", op)
		}
	}

	sum := 0
	for b := range boxes {
		for l := range boxes[b].l {
			sum += (b + 1) * (l + 1) * boxes[b].l[l].focal
		}
	}
	return fmt.Sprintf("%d", sum), nil
}

// 58375: too low

func day15Parse(rdr io.Reader) ([]string, error) {
	in, err := io.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	line := strings.Trim(string(in), "\n")
	seqList := strings.Split(line, ",")
	return seqList, nil
}

func (h *day15Hash) Write(b []byte) (int, error) {
	for _, c := range b {
		h.h += int(c)
		h.h = h.h * 17 % 256
	}
	return len(b), nil
}

func (h day15Hash) Sum() int {
	return h.h
}

func (b day15Box) Print(w io.Writer) {
	if len(b.l) > 0 {
		for l := range b.l {
			fmt.Fprintf(w, "[%s %d]", b.l[l].label, b.l[l].focal)
		}
	}
}
