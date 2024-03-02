package sqlbuilder

import (
	"io"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/delete"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/insert"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/sel"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/table"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/update"
)

type Formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}

type Builder struct {
	f        Formatter
	database string
}

func New(f Formatter) *Builder {
	return &Builder{
		f: f,
	}
}

func (b *Builder) SetDatabase(db string) *Builder {
	b.database = db
	return b
}

func (b *Builder) qualifiedTableName(table string) string {
	if b.database != `` {
		return b.database + `.` + table
	}
	return table
}

func (b *Builder) SelectFrom(tableExpr ast.IntoTableExpr) *sel.Builder {
	return sel.NewBuilder(b.f, tableExpr)
}

func (b *Builder) DeleteFromTable(table string) *delete.Builder {
	return delete.NewBuilder(b.f, b.qualifiedTableName(table))
}

func (b *Builder) UpdateTable(table string) *update.Builder {
	tableNode := ast.NewTableName(b.qualifiedTableName(table))
	return update.NewBuilder(b.f, tableNode)
}

func (b *Builder) InsertIntoTable(table string) *insert.Builder {
	return insert.NewBuilder(b.f, b.qualifiedTableName(table))
}

func (b *Builder) CreateTable(name string) *table.CreateBuilder {
	return table.NewCreateBuilder(b.f, name)
}

func (b *Builder) AlterTable(name string) *table.AlterBuilder {
	return table.NewAlterBuilder(b.f, name)
}
