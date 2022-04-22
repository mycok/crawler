package main

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	testcases := []struct {
		name     string
		cfg      config
		expected string
	}{
		{
			name:     "NoFilter",
			cfg:      config{root: "testdata", ext: "", list: true, size: 0},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n",
		},
		{
			name:     "FilterExtensionMatch",
			cfg:      config{root: "testdata", ext: ".log", list: true, size: 0},
			expected: "testdata/dir.log\n",
		},
		{
			name:     "FilterExtensionAndSizeMatch",
			cfg:      config{root: "testdata", ext: ".log", list: true, size: 10},
			expected: "testdata/dir.log\n",
		},
		{
			name:     "FilterExtensionButNoSizeMatch",
			cfg:      config{root: "testdata", ext: ".log", list: true, size: 20},
			expected: "",
		},
		{
			name:     "FilterExtensionWithNoMatch",
			cfg:      config{root: "testdata", ext: ".tar", list: true, size: 0},
			expected: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(&buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			output := buffer.String()

			if tc.expected != output {
				t.Errorf("Expected: %q, Got: %q", tc.expected, output)
			}
		})
	}
}
