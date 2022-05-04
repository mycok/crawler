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
	archive string
}

func main() {
	var (
		cfg config
		err error
		f   = os.Stdout
	)

	flag.StringVar(&cfg.root, "root", ".", "Root directory to start search")
	flag.StringVar(&cfg.ext, "ext", "", "File extension to filter out")
	flag.StringVar(&cfg.logfile, "log", "", "Filename to log all file deletions to")
	flag.StringVar(&cfg.archive, "archive", "", "Archive directory")
	flag.BoolVar(&cfg.del, "del", false, "Delete all files from the provided root directory")
	flag.BoolVar(&cfg.list, "ls", false, "List all files only")
	flag.Int64Var(&cfg.size, "size", 0, "Minimum file size in bytes")
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

		// Archive file / path to a specified archive directory and if successful continue to
		// other actions such delete if specified. Only return in case of an error.
		if cfg.archive != "" {
			if err := archiveFile(path, cfg.root, cfg.archive); err != nil {
				return err
			}
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
