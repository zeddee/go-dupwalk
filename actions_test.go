package main

import (
	"os"
	"testing"
)

func TestFilterOut(t *testing.T) {
	type testCase struct {
		name     string
		file     string
		ext      string
		exclude  string
		minSize  int64 // minimum file size in bytes
		expected bool
	}

	testCases := []testCase{
		{"FilterNoExtension", "test/data/dir.log", "", "", 0, false},
		{"FilterExtensionMatch", "test/data/dir.log", ".log", "", 0, false},
		{"FilterExtensionNoMatch", "test/data/dir.log", ".sh", "", 0, true},
		{"FilterExtensionSizeMatch", "test/data/dir.log", ".log", "", 10, false},
		{"FilterExtensionSizeNoMatch", "test/data/dir.log", ".log", "", 20, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			f := filterOut(tc.file, tc.ext, tc.exclude, tc.minSize, info)

			if f != tc.expected {
				t.Errorf("Expected '%t', got '%t' instead\n", tc.expected, f)
			}
		})
	}
}
