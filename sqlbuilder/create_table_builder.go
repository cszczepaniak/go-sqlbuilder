package sqlbuilder

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"

type createTableDialect interface {
	CreateTableStmt(name string) (string, error)
	CreateTableIfNotExistsStmt(name string) (string, error)
	ColumnStmt(c column.Column)
}

type CreateTableBuilder struct {
	name              string
	columns           []column.Column
	createIfNotExists bool
	ctd               createTableDialect
}

func createTable(d Dialect, name string) *CreateTableBuilder {
	return &CreateTableBuilder{
		name: name,
		ctd:  d,
	}
}

func (b *CreateTableBuilder) IfNotExists() *CreateTableBuilder {
	b.createIfNotExists = true
	return b
}

func (b *CreateTableBuilder) Columns(cs ...column.Column) *CreateTableBuilder {
	b.columns = append(b.columns, cs...)
	return b
}

func (b *CreateTableBuilder) Build() (string, error) {
	return ``, nil
}
