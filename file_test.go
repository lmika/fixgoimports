package main

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SrcFile", func() {
	It("should properly format imports", func() {
		f, err := NewGoFile(strings.NewReader(simpleSource))
		Expect(err).ToNot(HaveOccurred())

		f.SortImportsInPlace()

		sw := new(strings.Builder)
		err = f.Format(sw)
		Expect(err).ToNot(HaveOccurred())

		actualFormatted := sw.String()

		Expect(actualFormatted).To(Equal(expectedFormatted))
	})
})

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
