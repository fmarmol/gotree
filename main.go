package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var includeHidden bool

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

func walk(path string, prefix string, isLast bool) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	connector := "├── "
	if isLast {
		connector = "└── "
	}
	fmt.Println(prefix + connector + info.Name())

	if !info.IsDir() {
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	if !includeHidden {
		filtered := entries[:0]
		for _, e := range entries {
			if !isHidden(e.Name()) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	childPrefix := prefix + "│   "
	if isLast {
		childPrefix = prefix + "    "
	}

	for i, entry := range entries {
		walk(filepath.Join(path, entry.Name()), childPrefix, i == len(entries)-1)
	}
}

func main() {
	root := "."
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		if args[i] == "-i" {
			includeHidden = true
		} else {
			root = args[i]
		}
	}

	fmt.Println(root)

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !includeHidden {
		filtered := entries[:0]
		for _, e := range entries {
			if !isHidden(e.Name()) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	for i, entry := range entries {
		walk(filepath.Join(root, entry.Name()), "", i == len(entries)-1)
	}
}
