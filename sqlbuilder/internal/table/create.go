package table

import (
	"io"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
)

type Formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}

type CreateDialect interface {
	CreateTableStmt(name string) (string, error)
	CreateTableIfNotExistsStmt(name string) (string, error)
	ColumnStmt(c column.Column) (string, error)
	PrimaryKeyStmt(cs []string) (string, error)
}

type CreateBuilder struct {
	f Formatter

	name              string
	columns           []column.Builder
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

func (b *CreateBuilder) Columns(cs ...column.Builder) *CreateBuilder {
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
