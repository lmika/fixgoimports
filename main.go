package main

import (
	"bytes"
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var flagDryRun = flag.Bool("N", false, "Dry run")
	flag.Parse()

	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			if err := processFile(path, *flagDryRun); err != nil {
				log.Printf("%v: %v", info.Name(), err)
			}
		}
		return nil
	})
}

func processFile(filename string, dryRun bool) error {
	gf, err := NewGoFile(filename)
	if err != nil {
		return err
	}

	formattedFile := new(bytes.Buffer)
	if err := gf.Format(formattedFile); err != nil {
		return err
	}

	if !dryRun {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(f, formattedFile); err != nil {
			return err
		}
	}
	return nil
}