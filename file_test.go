package main

import (
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	actualFormatted := formatSource(t, simpleSource)

	if actualFormatted != expectedFormatted {
		t.Fatalf("actual formatted does not match expected: exp = [[%v]], actual = [[%v]]", expectedFormatted, actualFormatted)
	}
}

func formatSource(t *testing.T, src string) string {
	f, err := NewGoFile(strings.NewReader(src))
	if err != nil {
		t.Fatalf("expected no error while parsing but got %v", err)
	}
	sw := new(strings.Builder)
	if err := f.Format(sw); err != nil {
		t.Fatalf("expected no error while formatting but got %v", err)
	}
	return sw.String()
}

const simpleSource = `package main

import (
	"github.com/bla/fla"
	"regexp"
	"bitbucket.com/lmika/some-sample/foobar"
	"fmt"
	"net/http"
	"github.com/foo/bar"
)

func main() { }
`

const expectedFormatted = `package main

import (
	"fmt"
	"net/http"
	"regexp"

	"bitbucket.com/lmika/some-sample/foobar"
	"github.com/bla/fla"
	"github.com/foo/bar"
)

func main() { }
`
