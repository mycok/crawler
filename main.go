package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	root    string
	ext     string
	list    bool
	size    int64
	del     bool
	logfile string
	logger  io.Writer
}

func main() {
	var (
		cfg config
		err error
		f   = os.Stdout
	)

	flag.StringVar(&cfg.root, "r", ".", "Root directory to start search")
	flag.StringVar(&cfg.ext, "e", "", "File extension to filter out")
	flag.StringVar(&cfg.logfile, "log", "", "Log all file deletions to this file")
	flag.BoolVar(&cfg.del, "d", false, "File name / path to delete")
	flag.BoolVar(&cfg.list, "l", false, "List all files only")
	flag.Int64Var(&cfg.size, "s", 0, "Minimum file size in bytes")
	flag.Parse()

	if cfg.logfile != "" {
		f, err = os.OpenFile(cfg.logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			os.Exit(1)
		}

		defer f.Close()
	}

	cfg.logger = f

	if err = run(os.Stdout, cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

}

func run(out io.Writer, cfg config) error {
	var counter int64

	delLogger := log.New(cfg.logger, "FILE DELETED ON: ", log.LstdFlags)

	err := filepath.Walk(cfg.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, cfg.ext, cfg.size, info) {
			return nil
		}

		// If -l / list was explicitly set, don't do anything else.
		if cfg.list {
			return listFile(path, &counter, out)
		}

		// Delete matched files.
		if cfg.del {
			return delFile(path, &counter, delLogger)
		}

		// List is the default option if nothing else was set.
		return listFile(path, &counter, out)
	})

	displayMatchedCount(counter, cfg, out)

	return err
}
