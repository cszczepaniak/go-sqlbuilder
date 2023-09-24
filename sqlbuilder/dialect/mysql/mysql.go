package mysql

import (
	"fmt"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
)

func Now() string {
	return `NOW()`
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

	switch tc := c.(type) {
	case column.TinyIntColumn:
		sb.WriteString(`TINYINT`)
	case column.SmallIntColumn:
		sb.WriteString(`SMALLINT`)
	case column.IntColumn:
		sb.WriteString(`INT`)
	case column.BigIntColumn:
		sb.WriteString(`BIGINT`)
	case column.CharColumn:
		fmt.Fprintf(sb, `CHAR(%d)`, tc.Size)
	case column.VarCharColumn:
		fmt.Fprintf(sb, `VARCHAR(%d)`, tc.Size)
	case column.TextColumn:
		fmt.Fprintf(sb, `TEXT(%d)`, tc.Size)
	case column.TinyBlobColumn:
		sb.WriteString(`TINYBLOB`)
	case column.BlobColumn:
		sb.WriteString(`BLOB`)
	case column.MediumBlobColumn:
		sb.WriteString(`MEDIUMBLOB`)
	case column.LongBlobColumn:
		sb.WriteString(`LONGBLOB`)
	case column.DateTimeColumn:
		sb.WriteString(`DATETIME`)
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

	if column.AutoIncrement(c) {
		sb.WriteString(` AUTO_INCREMENT`)
	}

	return sb.String(), nil
}

func (m Dialect) PrimaryKeyStmt(cols []string) (string, error) {
	return `PRIMARY KEY(` + strings.Join(cols, `,`) + `)`, nil
}
