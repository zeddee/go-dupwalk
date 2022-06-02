package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type config struct {
	ext     string // File extensions to exclude from being listed
	exclude string // File paths and directories to exclude from being listed.
	// Uses filepath.Match, so uses fairly strict matching rules
	size       int64  // Bytes. Files larger than this are not listed
	list       bool   // Whether to print all found files to output
	outputFile string // Write list of duplicates out to a JSON file
	verbose    bool   // Prints to output each time an item is found
}

// run contains the main filepath.WalkDir implementation
func run(
	rootdir string, // Directory to start walking from
	output io.Writer, // Output destination to write to
	c config,
	dupList *[]duplicateFile,
) error {
	hist := &[]historyItem{}
	return filepath.WalkDir(
		rootdir,
		func(path string, dir os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			info, err := os.Stat(path)
			if err != nil {
				// if looking for duplicate files,
				// skip if cannot stat file instead of exiting
				fmt.Fprint(os.Stderr, err)
				return nil
			}

			switch {
			case filterOut(path, c.ext, c.exclude, c.size, info):
				// If path matches any filter rules,
				// continue
				return nil
			case c.list:
				// List files only
				return listFile(path, output)
			default:
				// Find all duplicate files
				*hist, *dupList, err = findDuplicates(path, output, c, hist, dupList)
				if err != nil {
					return err
				}
				return nil
			}
		})
}

func main() {
	outputFile := flag.String("out", "", "Write list of duplicates as JSON to this file.")
	root := flag.String("root", ".", "Start walking from this directory.")
	ext := flag.String("ext", "", "Filter out these file extensions.")
	exclude := flag.String("exclude", "", "Filter out files that match these filenames")
	size := flag.Int64("size", 0, "Minimum size of files to display.")
	list := flag.Bool("list", false, "List files only.")
	verbose := flag.Bool("v", false, "Prints to output each time an item is found")
	flag.Parse()

	c := config{
		outputFile: *outputFile,
		ext:        *ext,
		exclude:    *exclude,
		size:       *size,
		list:       *list,
		verbose:    *verbose,
	}

	dupList := &[]duplicateFile{}

	fmt.Fprint(os.Stderr, "Finding duplicates ...\n\n")

	if err := run(*root, os.Stderr, c, dupList); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	jsonOut, err := json.Marshal(dupList)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if c.outputFile != "" {
		fp, err := os.OpenFile(c.outputFile, os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		defer fp.Close()

		if _, err := fp.Write(jsonOut); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		fmt.Fprintln(os.Stdout, string(jsonOut))
	}
}
