package main

import (
	"flag"
	"os"
	"path/filepath"
)

func main() {
	var flagListDifferences = flag.Bool("l", false, "list files that differ from the formatted file")
	var flagWriteBack = flag.Bool("w", false, "write result to (source) file instead of stdout")
	flag.Parse()

	errs := &errorPresenter{}
	execContext := executionContext{
		writeBack:       *flagWriteBack,
		listFilesDiffer: *flagListDifferences,
		errorPresenter:  errs,
	}

	if flag.NArg() == 0 {
		if err := execContext.processStdin(); err != nil {
			errs.Errorf("stdin: %v", err)
		}
	} else {
		for _, arg := range flag.Args() {
			stat, err := os.Stat(arg)
			if err != nil {
				errs.Errorf("cannot stat %v: %v", arg, err)
				continue
			}

			if filepath.Ext(arg) == ".go" && !stat.IsDir() {
				err = execContext.processFile(arg)
			} else if stat.IsDir() {
				err = execContext.processDir(arg)
			} else {
				errs.Printf("ignoring non-go file %v", arg)
				err = nil
			}
			if err != nil {
				errs.Errorf("error %v: %v", arg, err)
			}
		}
	}
	errs.Exit()
}
