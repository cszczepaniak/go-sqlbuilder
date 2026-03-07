package integration

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cszczepaniak/gotest/assert"
)

// TestReadmeSnippetInSync fails if the README code block does not match the snippet extracted from
// readme_example_test.go. It generates the correct README. Commit the result.
func TestReadmeSnippetInSync(t *testing.T) {
	testPath := "readme_example_test.go"
	readmePath := "../README.md"

	want := readmeSnippetFromTestFile(t, testPath)
	got := readmeCodeBlock(t, readmePath)

	if got != want {
		t.Errorf("README.md code block is out of sync with readme_example_test.go. " +
			"This test will generate the correct contents; commit the result.",
		)
		updateReadme(t, want)
	}
}

const readmeExampleFunc = "TestReadmeExample"

func readmeSnippetFromTestFile(t *testing.T, testPath string) string {
	t.Helper()

	src, err := os.ReadFile(testPath)
	assert.NoError(t, err)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, testPath, src, parser.AllErrors)
	assert.NoError(t, err)

	tokFile := fset.File(f.Pos())

	b := &strings.Builder{}

	for _, d := range f.Decls {
		gen, ok := d.(*ast.GenDecl)
		if ok && gen.Tok == token.IMPORT {
			start := tokFile.Offset(d.Pos())
			end := tokFile.Offset(d.End())
			b.Write(src[start:end])
			b.WriteString("\n\n")
		}

		fn, ok := d.(*ast.FuncDecl)
		if ok && fn.Name.Name == readmeExampleFunc {
			start := tokFile.Offset(fn.Pos())
			end := tokFile.Offset(fn.End())
			b.Write(src[start:end])
		}
	}

	return b.String()
}

func readmeCodeBlock(t *testing.T, readmePath string) string {
	t.Helper()

	f, err := os.Open(readmePath)
	assert.NoError(t, err)
	defer f.Close()

	sc := bufio.NewScanner(f)
	var inBlock bool
	var b strings.Builder
	for sc.Scan() {
		line := sc.Text()
		if line == "```go" {
			inBlock = true
			continue
		}
		if inBlock {
			if line == "```" {
				break
			}
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}
	assert.NoError(t, sc.Err())

	return strings.TrimRight(b.String(), "\n")
}

func updateReadme(t *testing.T, snippet string) {
	readmePath := filepath.Join("..", "README.md")

	data, err := os.ReadFile(readmePath)
	assert.NoError(t, err)

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	start := strings.Index(content, "```go")
	if start == -1 {
		t.Fatal("could not find ```go block in README")
	}

	codeStart := start + len("```go")
	if len(content) > codeStart && content[codeStart] == '\n' {
		codeStart++
	}

	end := strings.Index(content[codeStart:], "\n```")
	if end == -1 {
		t.Fatal("could not find closing ``` in README")
	}

	end += codeStart
	newContent := content[:start] + "```go\n" + snippet + content[end:]

	if newContent == content {
		t.Log("README already in sync")
		return
	}

	assert.NoError(t, os.WriteFile(readmePath, []byte(newContent), 0o644))
	t.Log("Updated README.md")
}
