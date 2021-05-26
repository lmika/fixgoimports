package main

import (
	"github.com/pkg/errors"
	"go/scanner"
	"go/token"
	"strconv"
)

type srcFileReader struct {
	segments	[]segment
	srcBytes []byte
	fset	*token.FileSet
	scnr	*scanner.Scanner
	scanErr error

	nextPos token.Pos
	nextTok token.Token
	nextLit string
}

func newSrcFileReader(srcBytes []byte) *srcFileReader {
	var scnr scanner.Scanner

	fset := token.NewFileSet()
	file := fset.AddFile("test.go", 1, len(srcBytes))

	var srcFileReader = &srcFileReader{srcBytes: srcBytes, fset: fset}

	scnr.Init(file, srcBytes, func(pos token.Position, msg string) {
		srcFileReader.scanErr = errors.Errorf("scan error: %v", msg)
	}, scanner.ScanComments)

	srcFileReader.scnr = &scnr

	return srcFileReader
}

func (sr *srcFileReader) scanNext() {
	if sr.nextTok == token.EOF {
		return
	}
	sr.nextPos, sr.nextTok, sr.nextLit = sr.scnr.Scan()
}

func (sr *srcFileReader) consume(tok token.Token) error {
	if sr.nextTok != tok {
		return errors.Errorf("expected %v", tok)
	}
	sr.scanNext()
	return nil
}

func (sr *srcFileReader) consumeAny(tok token.Token){
	for sr.nextTok == tok {
		sr.scanNext()
	}
}

func (sr *srcFileReader) scanFile() error {
	for {
		before := sr.scanAndCollectUntil(token.IMPORT)
		sr.segments = append(sr.segments, codeSegment(before))

		if sr.nextTok == token.EOF {
			break
		}

		if err := sr.scanMultilineImport(); err != nil {
			return err
		}
	}

	return sr.scanErr
}

func (sr *srcFileReader) scanMultilineImport() error {
	if err := sr.consume(token.IMPORT); err != nil {
		return err
	}
	if err := sr.consume(token.LPAREN); err != nil {
		return err
	}

	var imports []Import
	for !sr.nextTokIsThisOrEOF(token.RPAREN) {
		nextImport, err := sr.scanImportExpression()
		if err != nil {
			return err
		}
		imports = append(imports, nextImport)

		sr.consumeAny(token.SEMICOLON)
	}

	if err := sr.consume(token.RPAREN); err != nil {
		return err
	}

	sr.segments = append(sr.segments, multiImportSegment{imports: imports})
	return nil
}

func (sr *srcFileReader) scanImportExpression() (imp Import, err error) {
	// TODO: scan comments
	var alias = ""
	if sr.nextTok == token.PERIOD {
		alias = "."
		sr.scanNext()
	} else if sr.nextTok == token.IDENT {
		alias = sr.nextLit
		sr.scanNext()
	}

	importLit := sr.nextLit
	if err := sr.consume(token.STRING); err != nil {
		return Import{}, err
	}

	importName, err := strconv.Unquote(importLit)
	if err != nil {
		return Import{}, err
	}

	return NewImport(importName, alias), nil
}

func (sr *srcFileReader) nextTokIsThisOrEOF(tok token.Token) bool {
	return sr.nextTok == tok || sr.nextTok == token.EOF
}


// scanAndCollectUntil collects all bytes from the source bytes up until (but not including) the token
// has been scanned, or the EOF has been reached
func (sr *srcFileReader) scanAndCollectUntil(tok token.Token) []byte {
	startOffset := sr.fset.Position(sr.nextPos).Offset
	for {
		sr.scanNext()
		if sr.nextTok == tok || sr.nextTok == token.EOF {
			break
		}
	}
	endOffset := sr.fset.Position(sr.nextPos).Offset
	return sr.srcBytes[startOffset:endOffset]
}


