package table

import (
	"io"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
)

type columnBuilder interface {
	Build() *ast.ColumnSpec
}

type Formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}

type CreateBuilder struct {
	f Formatter

	name              string
	columns           []columnBuilder
	createIfNotExists bool
}

func NewCreateBuilder(f Formatter, name string) *CreateBuilder {
	return &CreateBuilder{
		f:    f,
		name: name,
	}
}

func (b *CreateBuilder) IfNotExists() *CreateBuilder {
	b.createIfNotExists = true
	return b
}

func (b *CreateBuilder) Columns(cs ...columnBuilder) *CreateBuilder {
	b.columns = append(b.columns, cs...)
	return b
}

func (b *CreateBuilder) Build() (string, error) {
	ct := ast.NewCreateTable(b.name)
	if b.createIfNotExists {
		ct.CreateIfNotExists()
	}

	for _, col := range b.columns {
		ct.AddColumn(col.Build())
	}

	sb := &strings.Builder{}
	b.f.FormatNode(sb, ct)

	return sb.String(), nil
}
