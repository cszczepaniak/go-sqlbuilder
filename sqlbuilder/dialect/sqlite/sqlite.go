package sqlite

import (
	"fmt"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/conflict"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/params"
)

func Now() string {
	return `datetime('now')`
}

type Dialect struct{}

func (m Dialect) SelectStmt(table string, fields ...string) (string, error) {
	return `SELECT ` + strings.Join(fields, `,`) + ` FROM ` + table, nil
}

func (m Dialect) SelectForUpdateStmt(table string, fields ...string) (string, error) {
	// SQLite doesn't have select for update.
	return m.SelectStmt(table, fields...)
}

func (m Dialect) OrderBy(o filter.Order) (string, error) {
	s := `ORDER BY ` + o.Column + ` `
	switch o.Direction {
	case filter.Ascending:
		s += `ASC`
	case filter.Descending:
		s += `DESC`
	}
	return s, nil
}

func (m Dialect) DeleteStmt(table string) (string, error) {
	return `DELETE FROM ` + table, nil
}

func (m Dialect) UpdateStmt(table string, fields ...string) (string, error) {
	fieldList := &strings.Builder{}
	for i, f := range fields {
		fieldList.WriteString(f)
		fieldList.WriteString(`=?`)
		if i < len(fields)-1 {
			fieldList.WriteString(`,`)
		}
	}

	return `UPDATE ` + table + ` SET ` + fieldList.String(), nil
}

func (m Dialect) InsertStmt(table string, fields ...string) (string, error) {
	return `INSERT INTO ` + table + ` (` + strings.Join(fields, `,`) + `)`, nil
}

func (m Dialect) InsertIgnoreStmt(table string, fields ...string) (string, error) {
	return `INSERT OR IGNORE INTO ` + table + ` (` + strings.Join(fields, `,`) + `)`, nil
}

func (m Dialect) ValuesStmt(numRecords, numPerRecord int) (string, error) {
	return `VALUES ` + params.Groups(numRecords, numPerRecord), nil
}

func (m Dialect) Condition(f filter.Filter) (string, error) {
	c, err := m.condition(f)
	if err != nil {
		return ``, err
	}
	return `WHERE (` + c + `)`, nil
}

func (m Dialect) condition(f filter.Filter) (string, error) {
	if f == nil {
		return ``, nil
	}

	switch tf := f.(type) {
	case filter.EqualsFilter:
		return tf.Column + `=?`, nil
	case filter.NotEqualsFilter:
		return tf.Column + `!=?`, nil
	case filter.InFilter:
		return tf.Column + ` IN ` + params.Group(len(tf.Values)), nil
	case filter.GreaterFilter:
		return tf.Column + `>?`, nil
	case filter.GreaterOrEqualFilter:
		return tf.Column + `>=?`, nil
	case filter.LessFilter:
		return tf.Column + `<?`, nil
	case filter.LessOrEqualFilter:
		return tf.Column + `<=?`, nil
	case filter.AllFilter:
		return m.compositeCondition(tf.Filters, ` AND `)
	case filter.AnyFilter:
		return m.compositeCondition(tf.Filters, ` OR `)
	default:
		return ``, fmt.Errorf(`filter of type [%T] is not supported by MySQL`, f)
	}
}

func (m Dialect) compositeCondition(filters []filter.Filter, joinWith string) (string, error) {
	cs := make([]string, 0, len(filters))
	for _, ff := range filters {
		c, err := m.condition(ff)
		if err != nil {
			return ``, err
		}
		cs = append(cs, c)
	}
	return `(` + strings.Join(cs, joinWith) + `)`, nil
}

func (m Dialect) Limit() (string, error) {
	return `LIMIT ?`, nil
}

func (m Dialect) OnConflictStmt(key conflict.Key, conflicts ...conflict.Behavior) (string, error) {
	if len(conflicts) == 0 {
		return ``, nil
	}

	sb := &strings.Builder{}

	// Write the comma-delimited list of conflicting fields
	sb.WriteString(`ON CONFLICT (`)

	fields := key.Fields()
	for i, f := range fields {
		sb.WriteString(f)
		if i < len(fields)-1 {
			sb.WriteString(`,`)
		}
	}

	sb.WriteString(`)`)

	allIgnore := true
	for _, c := range conflicts {
		_, ok := c.(conflict.IgnoreBehavior)
		allIgnore = allIgnore && ok
	}

	if allIgnore {
		// Special case: if everything is ignored, sqlite supports DO NOTHING
		sb.WriteString(` DO NOTHING`)
		return sb.String(), nil
	}

	sb.WriteString(` DO UPDATE SET `)
	for i, c := range conflicts {
		sb.WriteString(c.Field())
		sb.WriteString(`=`)

		switch c.(type) {
		case conflict.IgnoreBehavior:
			sb.WriteString(c.Field())
		case conflict.OverwriteBehavior:
			sb.WriteString(`excluded.`)
			sb.WriteString(c.Field())
		}

		if i < len(conflicts)-1 {
			sb.WriteString(`,`)
		}
	}

	return sb.String(), nil
}

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
