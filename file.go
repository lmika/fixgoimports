package main

import (
	"io"
	"io/ioutil"
)

type GoFile struct {
	segments []segment
}

func NewGoFile(r io.Reader) (*GoFile, error) {
	bts, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	srcFileReader := newSrcFileReader(bts)
	if err := srcFileReader.scanFile(); err != nil {
		return nil, err
	}

	return &GoFile{segments: srcFileReader.segments}, nil
}

func (gf *GoFile) SortImportsInPlace() {
	for _, segment := range gf.segments {
		if mis, ok := segment.(multiImportSegment); ok {
			mis.imports.SortInPlace()
		}
	}
}

func (gf *GoFile) Format(w io.Writer) error {
	for _, segment := range gf.segments {
		if err := segment.format(w); err != nil {
			return err
		}
	}
	return nil
}
