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

func delFile(path string, counter *int64, out io.Writer) error {
	*counter++

	fmt.Fprint(out, "x")

	return os.Remove(path)
}

func displayMatchedCount(counter int64, cfg config, out io.Writer) {
	var str string

	subStr := "found"

	fmt.Fprintln(out)

	if counter == 1 {
		str = "%d file %s"
	} else {
		str = "%d files %s"
	}

	if cfg.del {
		subStr = "deleted"
	}


	fmt.Fprintln(out, fmt.Sprintf(str, counter, subStr))
}