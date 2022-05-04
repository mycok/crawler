package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}

func listFile(path string, counter *int64, out io.Writer) error {
	*counter++

	_, err := fmt.Fprintln(out, path)
	return err
}

func delFile(path string, counter *int64, delLogger *log.Logger) error {
	*counter++
	// Write to provided logger.
	delLogger.Println(path)

	return os.Remove(path)
}

func archiveFile(path, root, destDir string) error {
	// Retrieve destDir info.
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}

	// Check if the provided destDir is a directory.
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	// Retrieve the relative path of the provided path relative to the root.
	relativeDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	targetPath := filepath.Join(destDir, relativeDir, dest)

	// Create a destination directory.
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Open the output file at the target path for reading and writing.
	outputFile, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Open the file at the provided path for reading.
	inputFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	zw := gzip.NewWriter(outputFile)
	zw.Name = filepath.Base(path)

	if _, err := io.Copy(zw, inputFile); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}

	return outputFile.Close()
}

func displayMatchedCount(counter int64, cfg config, out io.Writer) {
	var mainStr string

	subStr := "found"

	fmt.Fprintln(out)

	if counter == 1 {
		mainStr = "%d file %s"
	} else {
		mainStr = "%d files %s"
	}

	if cfg.del {
		subStr = "deleted"
	}

	fmt.Fprintln(out, fmt.Sprintf(mainStr, counter, subStr))
}
