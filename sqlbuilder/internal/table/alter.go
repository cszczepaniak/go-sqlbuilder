package table

import (
	"context"
	"database/sql"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/dispatch"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type AlterBuilder struct {
	f Formatter

	name string

	columnsToAdd []columnBuilder
	indicesToAdd []indexBuilder
}

func NewAlterBuilder(f Formatter, name string) *AlterBuilder {
	return &AlterBuilder{
		f:    f,
		name: name,
	}
}

func (b *AlterBuilder) AddColumn(cs ...columnBuilder) *AlterBuilder {
	b.columnsToAdd = append(b.columnsToAdd, cs...)
	return b
}

func (b *AlterBuilder) AddIndex(is ...indexBuilder) *AlterBuilder {
	b.indicesToAdd = append(b.indicesToAdd, is...)
	return b
}

func (b *AlterBuilder) Build() (statement.Statement, error) {
	at := ast.NewAlterTable(b.name)

	for _, col := range b.columnsToAdd {
		at.AddColumn(col.Build())
	}

	for _, idx := range b.indicesToAdd {
		at.AddIndex(idx.Build())
	}

	sb := &strings.Builder{}
	b.f.FormatNode(sb, at)

	return statement.Statement{
		Stmt: sb.String(),
	}, nil
}

func (b *AlterBuilder) Exec(e dispatch.Execer) (sql.Result, error) {
	return dispatch.Exec(b, e)
}

func (b *AlterBuilder) ExecContext(ctx context.Context, e dispatch.ExecCtxer) (sql.Result, error) {
	return dispatch.ExecContext(ctx, b, e)
}
