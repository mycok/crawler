package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type config struct {
	root string
	ext  string
	list bool
	size int64
}

func main() {
	var cfg config

	flag.StringVar(&cfg.root, "r", ".", "Root directory to start search")
	flag.StringVar(&cfg.ext, "e", "", "File extension to filter out")
	flag.BoolVar(&cfg.list, "l", false, "List all files only")
	flag.Int64Var(&cfg.size, "s", 0, "Minimum file size in bytes")
	flag.Parse()

	if err := run(os.Stdout, cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

}

func run(out io.Writer, cfg config) error {
	return filepath.Walk(cfg.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, cfg.ext, cfg.size, info) {
			return nil
		}

		// If -l / list was explicitly set, don't do anything else.
		if cfg.list {
			return listFile(path, out)
		}

		// List is the default option if nothing else was set.
		return listFile(path, out)
	})
}
