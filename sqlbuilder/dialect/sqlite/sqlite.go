package sqlite

import (
	"fmt"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
)

func Now() string {
	return `datetime('now')`
}

type Dialect struct{}

func (m Dialect) CreateTableStmt(name string) (string, error) {
	return `CREATE TABLE ` + name, nil
}

func (m Dialect) CreateTableIfNotExistsStmt(name string) (string, error) {
	return `CREATE TABLE IF NOT EXISTS ` + name, nil
}

func (m Dialect) ColumnStmt(c column.Column) (string, error) {
	sb := &strings.Builder{}
	sb.WriteString(c.Name())
	sb.WriteString(` `)

	switch c.(type) {
	case column.TinyIntColumn,
		column.SmallIntColumn,
		column.IntColumn,
		column.BigIntColumn:
		sb.WriteString(`INTEGER`)
	case column.CharColumn,
		column.VarCharColumn,
		column.TextColumn:
		fmt.Fprintf(sb, `TEXT`)
	case column.TinyBlobColumn,
		column.BlobColumn,
		column.MediumBlobColumn,
		column.LongBlobColumn:
		sb.WriteString(`BLOB`)
	case column.DateTimeColumn:
		sb.WriteString(`NUMERIC`)
	}

	if n := c.Nullable(); n != nil && *n {
		sb.WriteString(` NULL`)
	} else if n != nil && !*n {
		sb.WriteString(` NOT NULL`)
	}

	if val, ok := c.Default(); ok {
		if column.IsText(c) {
			fmt.Fprintf(sb, ` DEFAULT %q`, val)
		} else {
			fmt.Fprintf(sb, ` DEFAULT %v`, val)
		}
	}

	return sb.String(), nil
}

func (m Dialect) PrimaryKeyStmt(cols []string) (string, error) {
	return `PRIMARY KEY(` + strings.Join(cols, `,`) + `)`, nil
}
