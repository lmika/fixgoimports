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
})

/*
func TestSrcFile(t *testing.T) {
	srcFile := newSrcFileReader([]byte(`package main

var somethingBefore int

import (
	"abc123"
)

var somethingAfter int`))

	err := srcFile.scanFile()
	if err != nil {
		t.Errorf("error while scanning file: %v", err)
	}
	t.Fatalf("print logs")
}
*/

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
