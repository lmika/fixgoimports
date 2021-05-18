package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	var flagWriteBack = flag.Bool("w", false, "write result to (source) file instead of stdout")
	flag.Parse()

	errs := errorPresenter{}

	if flag.NArg() == 0 {
		if err := processStdin(); err != nil {
			errs.Printf("stdin: %v", err)
		}
	} else {
		for _, arg := range flag.Args() {
			stat, err := os.Stat(arg)
			if err != nil {
				errs.Printf("cannot stat %v: %v", arg, err)
				continue
			}

			if filepath.Ext(arg) == ".go" && !stat.IsDir() {
				err = processFile(arg, *flagWriteBack)
			} else if stat.IsDir() {
				err = processDir(arg, *flagWriteBack)
			} else {
				errs.Warnf("ignoring non-go file %v", arg)
				err = nil
			}
			if err != nil {
				errs.Printf("error %v: %v", arg, err)
			}
		}
	}
	errs.Exit()
}

func processDir(dirName string, writeBack bool) error {
	return filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		filename := info.Name()

		// Skip files or directories starting with '.' or '_'
		if strings.HasPrefix(filename, ".") || strings.HasPrefix(filename, "_") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip vendor directories
		if filename == "vendor" && info.IsDir() {
			return filepath.SkipDir
		}

		if filepath.Ext(path) == ".go" {
			if err := processFile(path, writeBack); err != nil {
				fmt.Fprintf(os.Stderr, "%v: %v", path, err)
			}
		}
		return nil
	})
}

func processFile(filename string, writeBack bool) error {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "cannot read file")
	}

	gf, err := NewGoFile(bytes.NewReader(fileBytes))
	if err != nil {
		return errors.Wrap(err, "cannot read file")
	}

	formattedFile := new(bytes.Buffer)
	if err := gf.Format(formattedFile); err != nil {
		return errors.Wrap(err, "cannot format file")
	}

	if writeBack {
		f, err := os.Create(filename)
		if err != nil {
			return errors.Wrap(err, "cannot open file for writing")
		}
		defer f.Close()

		if _, err := io.Copy(f, formattedFile); err != nil {
			return errors.Wrap(err, "cannot write file")
		}
	} else {
		io.Copy(os.Stderr, formattedFile)
	}
	return nil
}

func processStdin() error {
	gf, err := NewGoFile(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "cannot read file")
	}

	if err := gf.Format(os.Stdout); err != nil {
		return errors.Wrap(err, "cannot format file")
	}
	return nil
}