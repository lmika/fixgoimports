package main

import "io"

type segment interface {
	format(w io.Writer) error
}

type codeSegment string

func (cs codeSegment) format(w io.Writer) (err error) {
	_, err = io.WriteString(w, string(cs))
	return
}

type singleImportSegment struct {
	theImport Import
}

func (mis singleImportSegment) format(w io.Writer) error {
	if _, err := io.WriteString(w, "import "); err != nil {
		return err
	}
	if err := mis.theImport.Format(w); err != nil {
		return err
	}

	return nil
}

type multiImportSegment struct {
	imports Imports
}

func (mis multiImportSegment) format(w io.Writer) error {
	if _, err := io.WriteString(w, "import (\n"); err != nil {
		return err
	}
	if err := mis.imports.Format(w); err != nil {
		return err
	}
	if _, err := io.WriteString(w, ")"); err != nil {
		return err
	}

	return nil
}
