package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var knownPrefixes = map[string]struct{}{
	"github.com":    {},
	"bitbucket.org": {},
	"pkg.go.dev":    {},
}

var blockImportName = regexp.MustCompile(`([.]|[a-zA-Z_][a-zA-Z0-9_]*)?\s*("[^"]*")`)

type Import struct {
	Name     string
	Alias string
	IsStdLib bool
}

var errBlank = errors.New("blank line")

func FromSourceLine(line string) (Import, error) {
	trimmedLine := strings.TrimSpace(line)
	if trimmedLine == "" {
		return Import{}, errBlank
	}

	submatches := blockImportName.FindStringSubmatch(trimmedLine)
	var quotedImport, alias string
	switch len(submatches) {
	case 2:
		quotedImport = submatches[1]
	case 3:
		alias = submatches[1]
		quotedImport = submatches[2]
	default:
		return Import{}, errors.Errorf("malformed line: %v", submatches)
	}

	unquotedImport, err := strconv.Unquote(quotedImport)
	if err != nil {
		return Import{}, err
	}

	return NewImport(unquotedImport, alias), nil
}

func NewImport(name, alias string) Import {
	firstSlash := strings.SplitN(name, "/", 2)
	if len(firstSlash) == 1 {
		return Import{Name: name, Alias: alias, IsStdLib: true}
	}

	if _, isKnownPrefix := knownPrefixes[firstSlash[0]]; isKnownPrefix {
		return Import{Name: name, Alias: alias, IsStdLib: false}
	}

	// Any dots in the first prefix likely to indicate a URL, meaning a non-stdlib name
	if strings.Count(firstSlash[0], ".") > 0 {
		return Import{Name: name, Alias: alias, IsStdLib: false}
	}

	return Import{Name: name, Alias: alias, IsStdLib: true}
}


type Imports	[]Import

func (imps Imports) SortInPlace() {
	sort.Slice(imps, func(i, j int) bool {
		if imps[i].IsStdLib && !imps[j].IsStdLib {
			return true
		} else if !imps[i].IsStdLib && imps[j].IsStdLib {
			return false
		} else {
			return imps[i].Name < imps[j].Name
		}
	})
}

func (imps Imports) Format(w io.Writer) error {
	lastStdLib := true
	for i, imp := range imps {
		if imp.IsStdLib != lastStdLib && i > 0 {
			if _, err := fmt.Fprintln(w, ""); err != nil {
				return err
			}
		}
		lastStdLib = imp.IsStdLib
		if imp.Alias != "" {
			if _, err := fmt.Fprintf(w, "\t%v %v\n", imp.Alias, strconv.Quote(imp.Name)); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, "\t%v\n", strconv.Quote(imp.Name)); err != nil {
				return err
			}
		}
	}
	return nil
}