package sqlbuilder

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/cszczepaniak/go-sqlbuilder/assert"
)

type exampleCode struct {
	Imports string
	Body    string
}

func extractExampleCode(t *testing.T, path string) exampleCode {
	f, err := os.Open(path)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, f.Close())
	})

	sc := bufio.NewScanner(f)
	importBuilder := strings.Builder{}
	bodyBuilder := strings.Builder{}

	inGoBlock := false
	inImports := false

	for sc.Scan() {
		if sc.Text() == "```go" {
			inGoBlock = true
			inImports = true
			continue
		}
		if !inGoBlock {
			continue
		}

		if sc.Text() == "```" {
			// We're at the end
			break
		}

		if inImports {
			if sc.Text() == `)` {
				fmt.Fprintln(&importBuilder, `"testing"`)
				fmt.Fprintln(&importBuilder, `"database/sql"`)
				inImports = false
			}
			fmt.Fprintln(&importBuilder, sc.Text())
		} else {
			fmt.Fprintln(&bodyBuilder, sc.Text())
		}
	}

	return exampleCode{
		Imports: importBuilder.String(),
		Body:    bodyBuilder.String(),
	}
}

func testExampleCode(t *testing.T, code exampleCode) {
	err := os.Mkdir(`doctest`, os.ModePerm)
	if !errors.Is(err, os.ErrExist) {
		assert.NoError(t, err)
	}

	f, err := os.OpenFile(filepath.Join(`doctest`, `docs_gen_test.go`), os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, f.Close())
		if os.Getenv(`CI`) == `true` {
			t.Log(`Removing temp files...`)
			assert.NoError(t, os.RemoveAll(`doctest`))
		}
	})

	tmpl, err := template.ParseFiles(`docs_test.tmpl`)
	assert.NoError(t, err)
	err = tmpl.Execute(f, code)
	assert.NoError(t, err)

	out, err := exec.Command(`go`, `test`, `./doctest`, `-v`, `-run`, `TestDocumentation`).CombinedOutput()
	assert.NoError(t, err)
	fmt.Println(string(out))
}

func TestReadmeExamples(t *testing.T) {
	readmePath := `../README.md`
	code := extractExampleCode(t, readmePath)
	testExampleCode(t, code)
}
