package sqlbuilder

import (
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
)

type createTableDialect interface {
	CreateTableStmt(name string) (string, error)
	CreateTableIfNotExistsStmt(name string) (string, error)
	ColumnStmt(c column.Column) (string, error)
	PrimaryKeyStmt(cs []string) (string, error)
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
	sb := &strings.Builder{}

	var createStmt string
	var err error
	if b.createIfNotExists {
		createStmt, err = b.ctd.CreateTableIfNotExistsStmt(b.name)
	} else {
		createStmt, err = b.ctd.CreateTableStmt(b.name)
	}
	if err != nil {
		return ``, nil
	}

	sb.WriteString(createStmt)
	sb.WriteString(`(`)

	pkCols := make([]string, 0)

	for _, c := range b.columns {
		if c.PrimaryKey() {
			pkCols = append(pkCols, c.Name())
		}

		cStr, err := b.ctd.ColumnStmt(c)
		if err != nil {
			return ``, err
		}
		sb.WriteString(cStr)
		sb.WriteString(`,`)
	}

	pkStmt, err := b.ctd.PrimaryKeyStmt(pkCols)
	if err != nil {
		return ``, err
	}
	sb.WriteString(pkStmt)

	sb.WriteString(`)`)
	return sb.String(), nil
}
