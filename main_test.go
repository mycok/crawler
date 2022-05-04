package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunFunctionality(t *testing.T) {
	t.Skip()
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

func TestDeleteFiles(t *testing.T) {
	t.Skip()

	testcases := []struct {
		name             string
		cfg              *config
		extNotToDelete   string
		numbOfDeleted    int
		numbOfNotDeleted int
		expected         string
	}{
		{
			name:             "DeleteWithNoExtMatch",
			cfg:              &config{ext: ".log", del: true},
			extNotToDelete:   ".gz",
			numbOfDeleted:    0,
			numbOfNotDeleted: 10,
			expected:         "0 files deleted",
		},
		{
			name:             "DeleteWithExtMatch",
			cfg:              &config{ext: ".log", del: true},
			extNotToDelete:   "",
			numbOfDeleted:    10,
			numbOfNotDeleted: 0,
			expected:         "\nXXXXXXXXXX\n10 files deleted",
		},
		{
			name:             "DeleteWithExtMixedUp",
			cfg:              &config{ext: ".log", del: true},
			extNotToDelete:   ".gz",
			numbOfDeleted:    5,
			numbOfNotDeleted: 5,
			expected:         "\nXXXXX\n5 files deleted",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)

			tempDir := createTempDir(t, map[string]int{
				tc.cfg.ext:        tc.numbOfDeleted,
				tc.extNotToDelete: tc.numbOfNotDeleted,
			})

			defer t.Cleanup(func() {
				os.RemoveAll(tempDir)
			})

			tc.cfg.root = tempDir
			tc.cfg.logger = &logBuffer

			if err := run(&buffer, *tc.cfg); err != nil {
				t.Fatal(err)
			}

			output := buffer.String()

			if tc.expected != output {
				t.Errorf("Expected: %q, Got: %q instead", tc.expected, output)
			}

			filesNotDeleted, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesNotDeleted) != tc.numbOfNotDeleted {
				t.Errorf("Expected: %d files not deleted, Got: %q instead", tc.numbOfNotDeleted, len(filesNotDeleted))
			}

			expectedLogLines := tc.numbOfDeleted + 1

			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))

			if len(lines) != expectedLogLines {
				t.Errorf("Expected %d log lines, Got: %d lines instead", expectedLogLines, len(lines))
			}
		})
	}
}

func TestArchiveFiles(t *testing.T) {
	testcases := []struct {
		name              string
		cfg               config
		extNotToArchive   string
		numbOfArchived    int
		numbOfNotArchived int
	}{
		{
			name:              "Extension to archive doesn't match",
			cfg:               config{ext: ".log"},
			extNotToArchive:   ".gz",
			numbOfArchived:    0,
			numbOfNotArchived: 10,
		},
		{
			name:              "Extension to archive matches",
			cfg:               config{ext: ".log"},
			extNotToArchive:   "",
			numbOfArchived:    10,
			numbOfNotArchived: 0,
		},
		{
			name:              "Archive mixed extensions",
			cfg:               config{ext: ".log"},
			extNotToArchive:   ".gz",
			numbOfArchived:    5,
			numbOfNotArchived: 5,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Buffer the archived output.
			var buffer bytes.Buffer

			// Create temp directories for both origin and archive directories.
			tempDir := createTempDir(t, map[string]int{
				tc.cfg.ext:         tc.numbOfArchived,
				tc.extNotToArchive: tc.numbOfNotArchived,
			})
			defer t.Cleanup(func() {
				os.RemoveAll(tempDir)
			})

			// Create a temp archive directory for writing archived files.
			tempArchiveDir := createTempDir(t, nil)
			defer t.Cleanup(func() {
				os.RemoveAll(tempArchiveDir)
			})

			tc.cfg.archive = tempArchiveDir
			tc.cfg.root = tempDir

			// Call the run function to perform file archiving.
			if err := run(&buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			// Create a search pattern to use to search for the archived files
			// from the dynamically created temp directory.
			searchPattern := filepath.Join(tempDir, fmt.Sprintf("*%s", tc.cfg.ext))

			// Perform a search for the archived files using the searchPattern.
			expectedFiles, err := filepath.Glob(searchPattern)
			if err != nil {
				t.Fatal(err)
			}

			expectedOutput := strings.Join(expectedFiles, "\n") + fmt.Sprintf("%d files found", tc.numbOfArchived)

			output := strings.TrimSpace(buffer.String())

			if expectedOutput != output {
				t.Errorf("Expected: %q, Got: %q instead", expectedOutput, output)
			}

			// Perform a check to validate the actual number of files archived.
			filesArchived, err := os.ReadDir(tempArchiveDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(filesArchived) != tc.numbOfArchived {
				t.Errorf("Expected: %d files archived, Got: %d files archived instead", tc.numbOfArchived, len(filesArchived))
			}
		})
	}
}

func createTempDir(t *testing.T, files map[string]int) (dirName string) {
	// Mark this fn as a test helper by calling t.Helper method.
	t.Helper()

	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDir, fname)

			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return tempDir
}
