package sqlbuilder

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/sel"

type Dialect interface {
	sel.Dialect
	deleteDialect
	updateDialect
	insertDialect
	createTableDialect
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

func (b *Builder) Delete(table string) *DeleteBuilder {
	return newDeleteBuilder(b.d, b.qualifiedTableName(table))
}

func (b *Builder) Update(table string) *UpdateBuilder {
	return newUpdateBuilder(b.d, b.qualifiedTableName(table))
}

func (b *Builder) Insert(table string) *InsertBuilder {
	return newInsertBuilder(b.d, b.qualifiedTableName(table))
}

func (b *Builder) CreateTable(name string) *CreateTableBuilder {
	return createTable(b.d, name)
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

func (b *TableBuilder) Delete() *DeleteBuilder {
	return newDeleteBuilder(b.b.d, b.table)
}

func (b *TableBuilder) Update() *UpdateBuilder {
	return newUpdateBuilder(b.b.d, b.table)
}

func (b *TableBuilder) Insert() *InsertBuilder {
	return newInsertBuilder(b.b.d, b.table)
}
