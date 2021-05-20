package main

import (
	"fmt"
	"os"
)

type errorPresenter struct {
	verbose bool
	errors  int
}

func (ep *errorPresenter) Verbosef(pattern string, args ...interface{}) {
	if ep.verbose {
		fmt.Fprintf(os.Stderr, pattern, args...)
		fmt.Fprintln(os.Stderr)
	}
}

func (ep *errorPresenter) Printf(pattern string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, pattern, args...)
	fmt.Fprintln(os.Stderr)
}

func (ep *errorPresenter) Errorf(pattern string, args ...interface{}) {
	ep.errors++
	fmt.Fprintf(os.Stderr, pattern, args...)
	fmt.Fprintln(os.Stderr)
}

func (ep *errorPresenter) Exit() {
	if ep.errors > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}
