package integration

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cszczepaniak/gotest/assert"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

// TestReadmeSnippetInSync fails if the README code block does not match the snippet extracted from
// readme_example_test.go. It generates the correct README. Commit the result.
func TestReadmeSnippetInSync(t *testing.T) {
	testPath := "readme_example_test.go"
	readmePath := "../README.md"
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		testPath = "integration/readme_example_test.go"
		readmePath = "README.md"
	}

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

func pkgPathsFromType(t types.Type, out map[string]bool) {
	switch t := t.(type) {
	case *types.Pointer:
		pkgPathsFromType(t.Elem(), out)
	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil {
			out[pkg.Path()] = true
		}
	}
}

func usedImportPaths(t *testing.T, fn *ast.FuncDecl, pkg *packages.Package) map[string]bool {
	used := make(map[string]bool)
	obj := pkg.Types.Scope().Lookup(readmeExampleFunc)
	if obj == nil {
		t.Fatal("test function not found")
	}

	sig := obj.Type().(*types.Signature)
	for i := 0; i < sig.Params().Len(); i++ {
		pkgPathsFromType(sig.Params().At(i).Type(), used)
	}

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		obj := pkg.TypesInfo.Uses[id]
		if obj == nil || obj.Pkg() == nil {
			return true
		}
		path := obj.Pkg().Path()
		if path != pkg.PkgPath {
			used[path] = true
		}
		return true
	})

	return used
}

func readmeSnippetFromTestFile(t *testing.T, testPath string) string {
	t.Helper()

	absPath, err := filepath.Abs(testPath)
	assert.NoError(t, err)

	dir := filepath.Dir(absPath)

	cfg := &packages.Config{
		Dir:   dir,
		Mode:  packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: true,
	}
	// Load the package that contains this test file (ensures we get the test variant).
	pkgs, err := packages.Load(cfg, "file="+absPath)
	assert.NoError(t, err)

	var pkg *packages.Package
	var obj types.Object
	var targetFile *ast.File
	var targetFunc *ast.FuncDecl

	for _, p := range pkgs {
		if p.Types == nil || p.TypesInfo == nil {
			continue
		}
		obj = p.Types.Scope().Lookup(readmeExampleFunc)
		if obj == nil {
			continue
		}
		if _, isFunc := obj.(*types.Func); !isFunc {
			continue
		}

		pkg = p
		break
	}

	if pkg == nil {
		t.Fatal("didn't find README example test")
	}

	insp := inspector.New(pkg.Syntax)
	fnIdentCursor, ok := insp.Root().FindByPos(obj.Pos(), obj.Pos())
	if !ok {
		t.Fatal("function not found")
	}
	targetFunc, ok = fnIdentCursor.Parent().Node().(*ast.FuncDecl)
	if !ok {
		t.Fatalf("function not a func decl: %T", fnIdentCursor.Node())
	}

	n := fnIdentCursor
	for n.Parent().Node() != nil {
		n = n.Parent()
	}
	targetFile, ok = n.Node().(*ast.File)
	if !ok {
		t.Fatal("file not a file")
	}

	tokFile := pkg.Fset.File(targetFunc.Pos())
	src, err := os.ReadFile(tokFile.Name())
	assert.NoError(t, err)

	used := usedImportPaths(t, targetFunc, pkg)

	snippet := &strings.Builder{}
	snippet.WriteString("import (\n")

	for _, d := range targetFile.Decls {
		gen, ok := d.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			continue
		}

		prevLine := -1
		for _, spec := range gen.Specs {
			imp, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			path := strings.Trim(imp.Path.Value, `"`)
			// Include if used in the function, or if it's a blank import (side-effect only).
			blankImport := imp.Name != nil && imp.Name.Name == "_"
			if !used[path] && !blankImport {
				continue
			}
			start := tokFile.Offset(imp.Pos())
			end := tokFile.Offset(imp.End())
			if start < 0 || end > len(src) {
				continue
			}

			// Preserve import groups
			line := pkg.Fset.Position(spec.Pos()).Line
			if prevLine != -1 && prevLine < line-1 {
				snippet.WriteString("\n")
			}

			snippet.WriteString("\t")
			snippet.WriteString(strings.TrimSpace(string(src[start:end])))
			snippet.WriteString("\n")

			prevLine = line
		}
		break
	}

	snippet.WriteString(")\n\n")

	body := targetFunc.Body

	bodyStart := tokFile.Offset(body.Lbrace) + 1
	firstTab := bytes.IndexByte(src[bodyStart:], '\t')
	bodyStart += firstTab

	bodyEnd := tokFile.Offset(body.Rbrace)

	for ln := range bytes.Lines(src[bodyStart:bodyEnd]) {
		snippet.Write(bytes.TrimPrefix(ln, []byte{'\t'}))
	}

	return strings.TrimRight(snippet.String(), "\n")
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
