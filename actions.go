package main

import (
	"fmt"
	"io"
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

func displayMatchedCount(counter int64, out io.Writer) {
	var str string

	fmt.Fprintln(out)

	if counter == 1 {
		str = "%d file found"
	} else {
		str = "%d files found"
	}

	fmt.Fprintln(out, fmt.Sprintf(str, counter))
}