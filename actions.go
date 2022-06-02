package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// filterOut applies takes a bunch rules
// to filter by, and returns true if it any of
// these rules match
func filterOut(
	path string, // Filepath currently being processed
	ext string, // Extensions to match
	exclude string, // Exclude pattern
	minSize int64, // Minimum file size
	info os.FileInfo, // Result of os.Stat(filepath)
) bool {
	// If filepath:
	// - is a directory, or
	// - is less than minimum filesize
	// return a match
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	// If ext matches the current filepath's extension, return a match
	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	// If an exclude pattern is used,
	// attempt to match the pattern to the current filepath.
	//
	// Match both full file path and file base+ext
	// for a more permissive matching policy,
	// because filepath.Match patterns require a strict match
	//
	// E.g. test/data/dir2/script.sh will only match the following patterns:
	// - `*/*/*/*.sh`
	// - `test/*/*/script.sh`
	// instead of being able to glob all directories like this:
	// - `*.sh`
	//
	// TODO:
	// - Get this to match dotfiles (e.g. `.DS_Store`)
	if exclude != "" {
		matched, err := filepath.Match(exclude, filepath.Clean(path))
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
		matched2, err := filepath.Match(
			exclude,
			fmt.Sprint(filepath.Base(path), filepath.Ext(path)),
		)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
		return matched || matched2
	}

	return false
}

// listFile takes a string and prints it out to a given io.Writer
func listFile(path string, output io.Writer) error {
	_, err := fmt.Fprintln(output, path)
	return err
}
