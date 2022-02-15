//line README.md:27
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/pgavlin/goldmark"
	"github.com/pgavlin/goldmark/ast"
	"github.com/pgavlin/goldmark/extension"
	goldmark_parser "github.com/pgavlin/goldmark/parser"
	"github.com/pgavlin/goldmark/text"
	"github.com/pgavlin/goldmark/util"
)
//line README.md:64
type lineIndex []int

func (index lineIndex) lineNumber(offset int) int {
	i := sort.Search(len(index), func(i int) bool {
		return index[i] > offset
	})
	return i + 1
}

func indexLines(f []byte) lineIndex {
	var index lineIndex
	for offset, b := range f {
		if b == '\n' {
			index = append(index, offset)
		}
	}
	return index
}
//line README.md:90
func convertFile(name string) error {
	contents, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	index := indexLines(contents)
//line README.md:102
	parser := goldmark.DefaultParser()
	parser.AddOptions(goldmark_parser.WithParagraphTransformers(
		util.Prioritized(extension.NewTableParagraphTransformer(), 200),
	))
	document := parser.Parse(text.NewReader(contents))
//line README.md:112
	var source bytes.Buffer
	ast.Walk(document, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		code, ok := n.(*ast.FencedCodeBlock)
		if !ok || !enter || string(code.Language(contents)) != "go" {
			return ast.WalkContinue, nil
		}

		lines := code.Lines()
		if lines.Len() == 0 {
			return ast.WalkContinue, nil
		}
//line README.md:129
		lineNumber := index.lineNumber(lines.At(0).Start)
		fmt.Fprintf(&source, "//line %v:%v\n", name, lineNumber)

		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			source.Write(line.Value(contents))
		}

		return ast.WalkContinue, nil
	})
//line README.md:145
	if source.Len() == 0 {
		return nil
	}
	return os.WriteFile(name+".go", source.Bytes(), 0600)
}
//line README.md:162
func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not read current directory: %v", err)
	}

	entries, err := os.ReadDir(wd)
	if err != nil {
		log.Fatalf("could not read current directory: %v", err)
	}
//line README.md:178
	for _, entry := range entries {
		name := entry.Name()
		ext := filepath.Ext(name)
		if ext != ".md" {
			continue
		}
		if err = convertFile(name); err != nil {
			log.Fatalf("could not convert file '%v': %v", name, err)
		}
	}
}
