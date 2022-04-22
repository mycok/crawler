package main

import (
	"os"
	"testing"
)

func TesTFilterOut(t *testing.T) {
	testcases := []struct {
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		{name: "FilterWithNoExtension", file: "testdata/dir.log", ext: "", minSize: 0, expected: false},
		{name: "FilterExtensionMatch", file: "testdata/dir.log", ext: ".log", minSize: 0, expected: false},
		{name: "FilterNoExtensionMath", file: "testdata/dir.log", ext: ".sh", minSize: 0, expected: true},
		{name: "FilterWithExtensionSizeMatch", file: "testdata/dir.log", ext: ".log", minSize: 10, expected: false},
		{name: "FilterWithExtensionNoSizeMatch", file: "testdata/dir.log", ext: ".log", minSize: 20, expected: true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal()
			}

			f := filterOut(tc.file, tc.ext, tc.minSize, info)
			if f != tc.expected {
				t.Errorf("Expected: %t', got: '%t' instead\n", tc.expected, f)
			}
		})
	}
}
