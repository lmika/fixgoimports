package main

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

var known3rdPartyPrefixes = map[string]struct{}{
	"github.com":    {},
	"bitbucket.org": {},
	"pkg.go.dev":    {},
}

type Import struct {
	Name     string
	Alias    string
	IsStdLib bool
}

func NewImport(name, alias string) Import {
	firstSlash := strings.SplitN(name, "/", 2)
	if len(firstSlash) == 1 {
		return Import{Name: name, Alias: alias, IsStdLib: true}
	}

	if _, isKnown3rdPartyPrefix := known3rdPartyPrefixes[firstSlash[0]]; isKnown3rdPartyPrefix {
		return Import{Name: name, Alias: alias, IsStdLib: false}
	}

	// Any dots in the first prefix likely to indicate a URL, meaning a non-stdlib name
	if strings.Count(firstSlash[0], ".") > 0 {
		return Import{Name: name, Alias: alias, IsStdLib: false}
	}

	return Import{Name: name, Alias: alias, IsStdLib: true}
}

func (imp Import) Format(w io.Writer) (err error) {
	if imp.Alias != "" {
		_, err = fmt.Fprintf(w, "%v %v", imp.Alias, strconv.Quote(imp.Name))
	} else {
		_, err = fmt.Fprint(w, strconv.Quote(imp.Name))
	}
	return
}

type Imports []Import

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

		if _, err := io.WriteString(w, "\t"); err != nil {
			return err
		}
		if err := imp.Format(w); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
	}
	return nil
}
