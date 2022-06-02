package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	type testCase struct {
		name     string
		root     string
		cfg      config
		expected string
	}

	testCases := []testCase{
		{
			name:     "NoFilter",
			root:     "test/data",
			cfg:      config{ext: "", size: 0, list: true},
			expected: "test/data/dir.log\ntest/data/dir2/script.sh\ntest/data/dir_duplicate.log\n"},
		{
			name:     "FilterExtensionMatch",
			root:     "test/data",
			cfg:      config{ext: ".log", size: 0, list: true},
			expected: "test/data/dir.log\ntest/data/dir_duplicate.log\n"},
		{
			name:     "FilterExtensionSizeMatch",
			root:     "test/data",
			cfg:      config{ext: ".log", size: 10, list: true},
			expected: "test/data/dir.log\ntest/data/dir_duplicate.log\n"},
		{
			name:     "FilterExtensionSizeNoMatch",
			root:     "test/data",
			cfg:      config{ext: ".log", size: 20, list: true},
			expected: ""},
		{
			name:     "FilterExtensionNoMatch",
			root:     "test/data",
			cfg:      config{ext: ".gz", size: 0, list: true},
			expected: ""},
		{
			name:     "FindDuplicates",
			root:     "test/data",
			cfg:      config{ext: "", size: 0, list: false, verbose: true},
			expected: "{\"Original\":\"test/data/dir.log\",\"Duplicate\":\"test/data/dir_duplicate.log\"}\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			if err := run(tc.root, &buffer, tc.cfg, &[]duplicateFile{}); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if tc.expected != res {
				t.Errorf("Expected %q, got %q instead\n", tc.expected, res)
			}
		})
	}
}
