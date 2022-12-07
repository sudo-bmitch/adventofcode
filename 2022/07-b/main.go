package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type dirEntry struct {
	name     string
	isDir    bool
	size     int
	parent   *dirEntry
	children map[string]*dirEntry
}

const (
	spaceTotal = 70000000
	spaceMin   = 30000000
)

func main() {
	cmd := "prompt"
	root := &dirEntry{
		name:     "/",
		isDir:    true,
		children: map[string]*dirEntry{},
	}
	var curDE *dirEntry
	in := bufio.NewScanner(os.Stdin)
	// build directory structure
	for in.Scan() {
		line := in.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "$ ") {
			cmd = "prompt"
		}
		switch cmd {
		case "prompt":
			cli := strings.Split(line, " ")
			if len(cli) < 2 {
				continue
			}
			if cli[1] == "cd" && len(cli) > 2 {
				if cli[2] == "/" {
					curDE = root
				} else if cli[2] == ".." && curDE != nil && curDE.parent != nil {
					curDE = curDE.parent
				} else if curDE != nil && curDE.children[cli[2]] != nil {
					curDE = curDE.children[cli[2]]
				} else {
					fmt.Fprintf(os.Stderr, "cd from %s to %s: not found", curDE.pwd(), cli[2])
					return
				}
			} else if cli[1] == "ls" {
				cmd = "ls"
			}
		case "ls":
			entry := strings.Split(line, " ")
			if len(entry) == 0 {
				continue
			} else if len(entry) != 2 {
				fmt.Fprintf(os.Stderr, "ls output unknown in %s: %s", curDE.pwd(), line)
				return
			}
			name := entry[1]
			newDE := &dirEntry{
				name:   name,
				parent: curDE,
			}
			if entry[0] == "dir" {
				newDE.isDir = true
				newDE.children = map[string]*dirEntry{}
			} else {
				size, err := strconv.Atoi(entry[0])
				if err != nil {
					fmt.Fprintf(os.Stderr, "ls size unknown in %s: %s: %v", curDE.pwd(), line, err)
					return
				}
				newDE.size = size
			}
			curDE.children[name] = newDE
		}
	}
	if err := in.Err(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "failed reading from stdin: %v", err)
		return
	}
	root.setDirSize()
	fmt.Printf("Tree:\n")
	root.printTree("")
	spaceNeeded := spaceMin - (spaceTotal - root.size)
	dir, size := root.smallestAbove(spaceNeeded)
	fmt.Printf("Result: needed %d, found %s with size %d\n", spaceNeeded, dir, size)
}

func (de *dirEntry) pwd() string {
	dir := "[nil]"
	curDE := de
	for curDE != nil {
		if dir == "[nil]" {
			dir = curDE.name
		} else if curDE.name == "/" {
			dir = "/" + dir
		} else {
			dir = curDE.name + "/" + dir
		}
		curDE = curDE.parent
	}
	return dir
}

func (de *dirEntry) setDirSize() {
	if de == nil {
		return
	}
	sum := 0
	for name, entry := range de.children {
		if entry.isDir {
			de.children[name].setDirSize()
		}
		sum += de.children[name].size
	}
	de.size = sum
}

func (de *dirEntry) printTree(prefix string) {
	if de == nil {
		return
	}
	deType := "file"
	if de.isDir {
		deType = "dir"
	}
	fmt.Printf("%s%s (%s, %d)\n", prefix, de.name, deType, de.size)
	for _, entry := range de.children {
		entry.printTree(prefix + "  ")
	}
}

func (de *dirEntry) smallestAbove(min int) (string, int) {
	if de == nil || !de.isDir {
		return "", 0
	}
	if de.size < min {
		return "", 0
	}
	bestName := de.pwd()
	bestSize := de.size
	for name, entry := range de.children {
		if entry.isDir {
			curName, curSize := de.children[name].smallestAbove(min)
			if curSize <= 0 {
				continue
			}
			if bestSize <= 0 || (curSize > 0 && curSize < bestSize) {
				bestName = curName
				bestSize = curSize
			}
		}
	}
	return bestName, bestSize
}
