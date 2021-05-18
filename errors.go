package main

import (
	"fmt"
	"os"
)

type errorPresenter struct {
	errors	int
}

func (ep *errorPresenter) Warnf(pattern string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, pattern, args...)
}

func (ep *errorPresenter) Printf(pattern string, args ...interface{}) {
	ep.errors++
	fmt.Fprintf(os.Stderr, pattern, args...)
}

func (ep *errorPresenter) Exit() {
	if ep.errors > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}