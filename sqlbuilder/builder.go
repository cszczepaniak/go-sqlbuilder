package sqlbuilder

import (
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/delete"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/insert"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/sel"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/table"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/update"
)

type Dialect interface {
	sel.Dialect
	delete.Dialect
	update.Dialect
	insert.Dialect
	table.CreateDialect
}

type Builder struct {
	d        Dialect
	database string
}

func New(d Dialect) *Builder {
	return &Builder{
		d: d,
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
	return sel.NewBuilder(b.d, target)
}

func (b *Builder) SelectFromTable(table string) *sel.Builder {
	target := sel.Table(b.qualifiedTableName(table))
	return b.SelectFrom(target)
}

func (b *Builder) Delete(table string) *delete.Builder {
	return delete.NewBuilder(b.d, b.qualifiedTableName(table))
}

func (b *Builder) Update(table string) *update.Builder {
	return update.NewBuilder(b.d, b.qualifiedTableName(table))
}

func (b *Builder) Insert(table string) *insert.Builder {
	return insert.NewBuilder(b.d, b.qualifiedTableName(table))
}

func (b *Builder) CreateTable(name string) *table.CreateBuilder {
	return table.NewCreateBuilder(b.d, name)
}

type TableBuilder struct {
	b     *Builder
	table string
}

func (b *Builder) ForTable(table string) *TableBuilder {
	return &TableBuilder{
		b:     b,
		table: table,
	}
}

func (b *TableBuilder) Select() *sel.Builder {
	return sel.NewBuilder(b.b.d, sel.Table(b.table))
}

func (b *TableBuilder) Delete() *delete.Builder {
	return delete.NewBuilder(b.b.d, b.table)
}

func (b *TableBuilder) Update() *update.Builder {
	return update.NewBuilder(b.b.d, b.table)
}

func (b *TableBuilder) Insert() *insert.Builder {
	return insert.NewBuilder(b.b.d, b.table)
}
