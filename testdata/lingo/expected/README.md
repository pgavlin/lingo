# lingo: literate programming with Go + Markdown

`lingo` is a simple tool for literate programming with Go and Markdown. `lingo` is
heavily inspired by [`tango`](https://github.com/pnkfelix/tango), a similar tool designed
for literate programming with Rust and Markdown.

When run, lingo` will extract Go source code from fenced code blocks in each Markdown
file in the current directory. Markdown files must use the `.md` extension, and code
will only be extracted from fenced code blocks with the language `go`. Each Markdown
file `some-file.md` that contains Go code will be converted into a file `some-file.md.go`.

To author a program with `lingo`, simply write your program as fenced code blocks in
Markdown files, then add a `.go` file in the same directory with a `//go:generate lingo`
directive preceding its package name.

This file is the source for `lingo` itself; let's break it down!

## Preamble

As usual, we start our program with a package clause followed by our import declarations.
Because we're going to be working with Markdown, our only imports outside the standard
library are from a [fork](https://github.com/pgavlin/goldmark) of the
[Goldmark](https://github.com/yuin/goldmark) Markdown parser. We'll be using that package
to parse Markdown into an AST that we'll then use as a basis for source code extraction.

```go
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
```

## Source Position Mapping

Because we're essentially generating source code, we'd like the extracted source code to
retain its original source positions. This allows downstream tools to reference positions
in the Markdown rather than positions in the extracted code. Go gives us the ability to
propagate this information through the use of [line directives](https://pkg.go.dev/cmd/compile#hdr-Compiler_Directives).

The only position information we need is the line number itself, as we'll be emitting
directives of the form `//line filename:line`. Unfortunately, Goldmark does not track
line information in its AST! It does, however, track the byte offset of each block of
text, including the contents of code blocks. We can determine the line number of a code
block ourselves by first building a byte offset to line number index from the Markdown
source. This index is a simple list of integers, where each entry tracks `E_i` is the
byte offset of the end of line `i`. With this structure, we can determine the number of
the line that contains a particular offset `o` by searching for the smallest index `i`
where `E_i > o`; the 1-indexed line number containing `o` is then `i + 1`.

```go
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
```

## Extracting Go Source from Markdown

With our line index implemented, converting each file is straightforward. First, we read
in the source code and build our line index:

```go
func convertFile(name string) error {
	contents, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	index := indexLines(contents)
```

Next, we parse the Markdown:

```go
	parser := goldmark.DefaultParser()
	parser.AddOptions(goldmark_parser.WithParagraphTransformers(
		util.Prioritized(extension.NewTableParagraphTransformer(), 200),
	))
	document := parser.Parse(text.NewReader(contents))
```

Then, we walk the parsed AST, looking for fenced code blocks with the language `go`:

```go
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
```

When we find a suitable code block, we determine its line number, then emit a line
directive followed by the contents of the code block into our output:

```go
		lineNumber := index.lineNumber(lines.At(0).Start)
		fmt.Fprintf(&source, "//line %v:%v\n", name, lineNumber)

		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			source.Write(line.Value(contents))
		}

		return ast.WalkContinue, nil
	})
```

Finally, we emit the collected source into an output file and return. If the walk did not
extract any source code, we do not emit an output file.

```go
	if source.Len() == 0 {
		return nil
	}
	return os.WriteFile(name+".go", source.Bytes(), 0600)
}
```

## The Entry Point

The only thing left to do now is to implement `lingo`'s entry point. The entry point is
responsible for finding the Markdown files the tool will convert and driving their
conversion using [`convertFile`](#extracting-go-source-from-markdown).

`lingo` operates on Markdown files in the current directory, so we begin by fetching the
path of the current directory and listing its contents:

```go
func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not read current directory: %v", err)
	}

	entries, err := os.ReadDir(wd)
	if err != nil {
		log.Fatalf("could not read current directory: %v", err)
	}
```

Then, we iterate the directory's contents and attempt to convert each `.md` file to a
`.md.go` file.

```go
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
```

And we're done!
