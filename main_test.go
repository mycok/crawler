package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n\n2 files found\n",
		},
		{
			name:     "FilterExtensionMatch",
			cfg:      config{root: "testdata", ext: ".log", list: true, size: 0},
			expected: "testdata/dir.log\n\n1 file found\n",
		},
		{
			name:     "FilterExtensionAndSizeMatch",
			cfg:      config{root: "testdata", ext: ".log", list: true, size: 10},
			expected: "testdata/dir.log\n\n1 file found\n",
		},
		{
			name:     "FilterExtensionButNoSizeMatch",
			cfg:      config{root: "testdata", ext: ".log", list: true, size: 20},
			expected: "\n0 files found\n",
		},
		{
			name:     "FilterExtensionWithNoMatch",
			cfg:      config{root: "testdata", ext: ".tar", list: true, size: 0},
			expected: "\n0 files found\n",
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
				t.Errorf("Expected: %q, Got: %q instead", tc.expected, output)
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	testcases := []struct {
		name        string
		cfg         *config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{
			name:        "DeleteWithNoExtMatch",
			cfg:         &config{ext: ".log", del: true},
			extNoDelete: ".gz",
			nDelete:     0,
			nNoDelete:   10,
			expected:    "0 files deleted",
		},
		{
			name:        "DeleteWithExtMatch",
			cfg:         &config{ext: ".log", del: true},
			extNoDelete: "",
			nDelete:     10,
			nNoDelete:   0,
			expected:    "\nXXXXXXXXXX\n10 files deleted",
		},
		{
			name:        "DeleteWithExtMixedUp",
			cfg:         &config{ext: ".log", del: true},
			extNoDelete: ".gz",
			nDelete:     5,
			nNoDelete:   5,
			expected:    "\nXXXXX\n5 files deleted",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)

			tempDirName := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})

			defer t.Cleanup(func() {
				os.RemoveAll(tempDirName)
			})

			tc.cfg.root = tempDirName
			tc.cfg.logger = &logBuffer

			if err := run(&buffer, *tc.cfg); err != nil {
				t.Fatal(err)
			}

			output := buffer.String()

			if tc.expected != output {
				t.Errorf("Expected: %q, Got: %q instead", tc.expected, output)
			}

			filesNotDeleted, err := os.ReadDir(tempDirName)
			if err != nil {
				t.Error(err)
			}

			if len(filesNotDeleted) != tc.nNoDelete {
				t.Errorf("Expected: %d files not deleted, Got: %q instead", tc.nNoDelete, len(filesNotDeleted))
			}

			expLogLines := tc.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expLogLines {
				t.Errorf("Expected %d log lines, Got: %d lines instead", expLogLines, len(lines))
			}
		})
	}
}

func createTempDir(t *testing.T, files map[string]int) (dirName string) {
	// Mark this fn as a test helper by calling t.Helper method.
	t.Helper()

	tempDir, err := os.CreateTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	tempDirName := tempDir.Name()

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDirName, fname)

			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return tempDirName
}
