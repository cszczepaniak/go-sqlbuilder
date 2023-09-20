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

type Dialect interface {
	delete.Dialect
	update.Dialect
	insert.Dialect
	table.CreateDialect
}

type Formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}

type Builder struct {
	d        Dialect
	f        Formatter
	database string
}

func New(d Dialect, f Formatter) *Builder {
	return &Builder{
		d: d,
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

func (b *Builder) SelectFrom(target sel.Target) *sel.Builder {
	return sel.NewBuilder(b.f, target)
}

func (b *Builder) SelectFromTable(table string) *sel.Builder {
	target := sel.Table(b.qualifiedTableName(table))
	return b.SelectFrom(target)
}

func (b *Builder) DeleteFromTable(table string) *delete.Builder {
	return delete.NewBuilder(b.d, b.qualifiedTableName(table))
}

func (b *Builder) UpdateTable(table string) *update.Builder {
	return update.NewBuilder(b.d, b.qualifiedTableName(table))
}

func (b *Builder) InsertIntoTable(table string) *insert.Builder {
	return insert.NewBuilder(b.d, b.qualifiedTableName(table))
}

func (b *Builder) CreateTable(name string) *table.CreateBuilder {
	return table.NewCreateBuilder(b.d, name)
}
