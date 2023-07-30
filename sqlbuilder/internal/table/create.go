package table

import (
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
)

type CreateDialect interface {
	CreateTableStmt(name string) (string, error)
	CreateTableIfNotExistsStmt(name string) (string, error)
	ColumnStmt(c column.Column) (string, error)
	PrimaryKeyStmt(cs []string) (string, error)
}

type CreateBuilder struct {
	name              string
	columns           []column.Column
	createIfNotExists bool
	ctd               CreateDialect
}

func NewCreateBuilder(d CreateDialect, name string) *CreateBuilder {
	return &CreateBuilder{
		name: name,
		ctd:  d,
	}
}

func (b *CreateBuilder) IfNotExists() *CreateBuilder {
	b.createIfNotExists = true
	return b
}

func (b *CreateBuilder) Columns(cs ...column.Column) *CreateBuilder {
	b.columns = append(b.columns, cs...)
	return b
}

func (b *CreateBuilder) Build() (string, error) {
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
