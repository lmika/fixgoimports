package main

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

type GoFile struct {
	Before string		// Lines before and including 'import ('
	Imports Imports // Lines within 'import (' and ')'
	After string
}

func NewGoFile(filename string) (*GoFile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	state := 0
	before := new(strings.Builder)
	after := new(strings.Builder)
	importLines := make(Imports, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		trimmedLine := strings.TrimSpace(scanner.Text())

		switch state {
		case scanStateBefore:
			before.WriteString(scanner.Text())
			before.WriteRune('\n')
			if trimmedLine == "import (" {
				state = scanStateWithin
			}
		case scanStateWithin:
			if trimmedLine == ")" {
				after.WriteString(")\n")
				state = scanStateAfter
			} else if trimmedLine != "" {
				importName, err := FromSourceLine(trimmedLine)
				if err != nil {
					return nil, errors.Wrapf(err, "malformed import line: %v", trimmedLine)
				}
				importLines = append(importLines, importName)
			}
		case scanStateAfter:
			if trimmedLine == "import (" {
				return nil, errors.New("found two import blocks")
			}
			after.WriteString(scanner.Text())
			after.WriteRune('\n')
		}
	}

	importLines.SortInPlace()

	return &GoFile{
		Before: before.String(),
		Imports: importLines,
		After: after.String(),
	}, nil
}

func (gf *GoFile) Format(w io.Writer) error {
	if _, err := io.WriteString(w, gf.Before); err != nil {
		return err
	}
	if err := gf.Imports.Format(w); err != nil {
		return err
	}
	if _, err := io.WriteString(w, gf.After); err != nil {
		return err
	}
	return nil
}

const (
	scanStateBefore int = iota
	scanStateWithin
	scanStateAfter
)