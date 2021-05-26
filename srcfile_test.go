package main

import (
	"bufio"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("SrcFile", func() {
	It("properly parse the file", func() {
		srcFile := newSrcFileReader([]byte(deindent(`
			var somethingBefore int
			import (
				"abc123"
			)
			var somethingAfter int
		`)))

		Expect(srcFile.scanFile()).ToNot(HaveOccurred())

		Expect(srcFile.segments).To(HaveLen(3))
		Expect(srcFile.segments[0]).To(Equal(codeSegment("var somethingBefore int\n")))
		Expect(srcFile.segments[1]).To(Equal(multiImportSegment{
			imports: []Import{
				{Name: "abc123", IsStdLib: true},
			},
		}))
		Expect(srcFile.segments[2]).To(Equal(codeSegment("\nvar somethingAfter int\n")))
	})

	It("should return an error if a scan error occurred", func() {
		srcFile := newSrcFileReader([]byte(deindent(`
			# some invalid character
		`)))

		Expect(srcFile.scanFile()).To(HaveOccurred())
	})

	It("properly parse the file with multiple imports", func() {
		srcFile := newSrcFileReader([]byte(deindent(`
			var somethingBefore int
			// and comments and such
			import (
				"strings"

				"example.com/something/else"

				. "bitbucket.com/import/by/default"
				cool_err "github.com/pkg/errors"
			)
			var somethingAfter int

			func main() { }
		`)))

		Expect(srcFile.scanFile()).ToNot(HaveOccurred())

		Expect(srcFile.segments).To(HaveLen(3))
		Expect(srcFile.segments[0]).To(Equal(codeSegment("var somethingBefore int\n// and comments and such\n")))
		Expect(srcFile.segments[1]).To(Equal(multiImportSegment{
			imports: []Import{
				{Name: "strings", IsStdLib: true},
				{Name: "example.com/something/else", IsStdLib: false},
				{Name: "bitbucket.com/import/by/default", IsStdLib: false, Alias: "."},
				{Name: "github.com/pkg/errors", IsStdLib: false, Alias: "cool_err"},
			},
		}))
		Expect(srcFile.segments[2]).To(Equal(codeSegment("\nvar somethingAfter int\n\nfunc main() { }\n")))
	})

	It("support multiple import statements", func() {
		srcFile := newSrcFileReader([]byte(deindent(`
			package main

			import (
				"strings"
			)
			import (
				"example.com/something/else"
			)

			func main() { }
		`)))

		Expect(srcFile.scanFile()).ToNot(HaveOccurred())

		Expect(srcFile.segments).To(HaveLen(5))
		Expect(srcFile.segments[0]).To(Equal(codeSegment("package main\n\n")))
		Expect(srcFile.segments[1]).To(Equal(multiImportSegment{
			imports: []Import{{Name: "strings", IsStdLib: true}},
		}))
		Expect(srcFile.segments[2]).To(Equal(codeSegment("\n")))
		Expect(srcFile.segments[3]).To(Equal(multiImportSegment{
			imports: []Import{{Name: "example.com/something/else", IsStdLib: false}},
		}))
		Expect(srcFile.segments[4]).To(Equal(codeSegment("\n\nfunc main() { }\n")))
	})

	It("should not parse strings as if they were imports", func() {
		srcFile := newSrcFileReader([]byte(deindent(`
			package main

			import (
				"strings"
			)

			var fancyString = ` + "`" + `
				import (
					"example.com/something/else"
				)
			` + "`" + `

			func main() { }
		`)))

		Expect(srcFile.scanFile()).ToNot(HaveOccurred())

		Expect(srcFile.segments).To(HaveLen(3))
		Expect(srcFile.segments[0]).To(Equal(codeSegment("package main\n\n")))
		Expect(srcFile.segments[1]).To(Equal(multiImportSegment{
			imports: []Import{{Name: "strings", IsStdLib: true}},
		}))
		Expect(srcFile.segments[2]).To(BeAssignableToTypeOf(codeSegment("")))
	})

})

// deindent takes indentation queues from the first line of the string, and removes it from all
// successive lines.  If the first line is blank, the line is also stripped.
func deindent(str string) string {
	scnr := bufio.NewScanner(strings.NewReader(str))
	outStr := new(strings.Builder)

	firstLine := true
	indentPos := 0
	for scnr.Scan() {
		line := scnr.Text()

		if firstLine && line == "" {
			continue
		} else if firstLine {
			indentPos = strings.IndexFunc(line, func(r rune) bool { return r != ' ' && r != '\t' })
			if indentPos < 0 {
				indentPos = 0
			}
			firstLine = false
		} else {
			outStr.WriteRune('\n')
		}

		firstNonSpace := strings.IndexFunc(line, func(r rune) bool { return r != ' ' && r != '\t' })
		if firstNonSpace == -1 {
			line = ""
		} else if firstNonSpace > indentPos {
			line = line[indentPos:]
		} else {
			line = line[firstNonSpace:]
		}
		outStr.WriteString(line)
	}
	return outStr.String()
}
