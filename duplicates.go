package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
)

// historyItems is used to compose a
// list of files that have already been read
// when listing duplicate files
type historyItem struct {
	path   string
	sha256 fileHash
}

// duplicateFile is used to contain and marshal
// to JSON a set of filepaths that describes
// a duplicate file
type duplicateFile struct {
	Original  string
	Duplicate string
}

// fileHash is an alias for a sha256 hash
type fileHash [32]byte // sha256 hash

// hashFile reads the first 5MB/5120 bytes
// of a file and produces a filehash,
// which is then used by inHistory to
// check if a file is a duplicate of a
// previously read file
func hashFile(path string) fileHash {
	fp, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	buffer := make([]byte, 5120) // read only the first 5120 bytes
	if _, err := fp.Read(buffer); err != nil && err != io.EOF {
		fmt.Fprint(os.Stderr, fmt.Errorf("read error: %s", err))
		os.Exit(1)
	}

	if err := fp.Close(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	return sha256.Sum256(buffer)
}

// inHistory checks if the currently read file
// has the same hash as any previously read file
func inHistory(thisFile historyItem, curHist []historyItem) (duplicateFile, bool) {
	for _, item := range curHist {
		if reflect.DeepEqual(thisFile.sha256, item.sha256) {
			return duplicateFile{
				Original:  item.path,
				Duplicate: thisFile.path,
			}, true
		}
	}
	return duplicateFile{}, false
}

// findDuplicates takes the most up-to-date list of previously-files
// and attempts to match the file at the currently path
// to items in that list.
//
// If it finds a positive match, findDuplicates
// prints a duplicateFile JSON object
// to output.
func findDuplicates(
	path string,
	output io.Writer,
	c config,
	currentHistory *[]historyItem,
	dupList *[]duplicateFile,
) ([]historyItem, []duplicateFile, error) {
	thisFile := historyItem{
		path:   path,
		sha256: hashFile(path),
	}

	if len(*currentHistory) == 0 {
		*currentHistory = append(*currentHistory, thisFile)
		return *currentHistory, *dupList, nil // don't write to output unless it is a duplicate file
	}

	dup, isDup := inHistory(thisFile, *currentHistory)

	if isDup {
		jsonOut, err := json.Marshal(dup)
		if err != nil {
			return *currentHistory, *dupList, err
		}
		*dupList = append(*dupList, dup)
		if c.verbose {
			fmt.Fprintf(output, "%s\n", string(jsonOut))
		}
		return *currentHistory, *dupList, nil
	}

	*currentHistory = append(*currentHistory, thisFile)
	return *currentHistory, *dupList, nil // don't write to output unless it is a duplicate file
}
