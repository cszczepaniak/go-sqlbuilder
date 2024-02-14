package formatter

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
	"github.com/stretchr/testify/assert"
)

type formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}

type formatTestCase struct {
	f    formatter
	node ast.Node
	exp  string
}

func newFormatTestCase(f formatter, node ast.Node, exp string) formatTestCase {
	return formatTestCase{
		f:    f,
		node: node,
		exp:  exp,
	}
}

func assertFormatting(t *testing.T, cases ...formatTestCase) {
	t.Helper()

	for _, c := range cases {
		t.Run(fmt.Sprintf("%T", c.f), func(t *testing.T) {
			sb := &strings.Builder{}
			c.f.FormatNode(sb, c.node)
			assert.Equal(t, c.exp, sb.String())
		})
	}
}

func assertAllFormatting(t *testing.T, node ast.Node, exp string) {
	t.Helper()

	assertFormatting(
		t,
		newFormatTestCase(Mysql{}, node, exp),
		newFormatTestCase(Sqlite{}, node, exp),
	)
}

func TestTableAlias(t *testing.T) {
	node := &ast.TableAlias{
		Alias: &ast.Alias{
			ForExpr: ast.NewIdentifier("foo"),
			As:      ast.NewIdentifier("bar"),
		},
	}

	assertAllFormatting(t, node, `foo AS bar`)
}

func TestSelector(t *testing.T) {
	node := &ast.Selector{
		SelectFrom: ast.NewIdentifier("foo"),
		FieldName:  ast.NewIdentifier("bar"),
	}

	assertAllFormatting(t, node, "foo.bar")
}
