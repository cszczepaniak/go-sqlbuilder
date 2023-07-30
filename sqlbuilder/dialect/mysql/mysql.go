package mysql

import (
	"fmt"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/conflict"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/functions"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/expr"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/params"
)

func Now() string {
	return `NOW()`
}

type Dialect struct{}

func (m Dialect) ResolveExpr(ex expr.Expr) (string, error) {
	switch te := ex.(type) {
	case expr.Column:
		if te.IsQualified() {
			return te.Database + `.` + te.Name, nil
		} else {
			return te.Name, nil
		}

	case functions.Count:
		if te.All() {
			return `COUNT(*)`, nil
		}
		c := `COUNT(`
		if te.Distinct {
			c += `DISTINCT `
		}
		c += te.Field
		c += `)`
		return c, nil
	}

	return ``, fmt.Errorf(`unsupported expression type: %T`, ex)
}

func (m Dialect) SelectStmt(table string, fields ...expr.Expr) (string, error) {
	return m.selectStmt(table, fields...)
}

func (m Dialect) SelectForUpdateStmt(table string, fields ...expr.Expr) (string, error) {
	stmt, err := m.selectStmt(table, fields...)
	if err != nil {
		return ``, err
	}
	return stmt + ` FOR UPDATE`, nil
}

func (m Dialect) selectStmt(table string, fields ...expr.Expr) (string, error) {
	resolved := make([]string, 0, len(fields))
	for _, f := range fields {
		r, err := m.ResolveExpr(f)
		if err != nil {
			return ``, err
		}

		resolved = append(resolved, r)
	}
	return `SELECT ` + strings.Join(resolved, `,`) + ` FROM ` + table, nil
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
	return `INSERT IGNORE INTO ` + table + ` (` + strings.Join(fields, `,`) + `)`, nil
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

func (m Dialect) OnConflictStmt(_ conflict.Key, conflicts ...conflict.Behavior) (string, error) {
	if len(conflicts) == 0 {
		return ``, nil
	}

	sb := &strings.Builder{}
	sb.WriteString(`ON DUPLICATE KEY UPDATE `)
	for i, c := range conflicts {
		sb.WriteString(c.Field())
		sb.WriteString(`=`)

		switch c.(type) {
		case conflict.IgnoreBehavior:
			sb.WriteString(c.Field())
		case conflict.OverwriteBehavior:
			sb.WriteString(`VALUES(`)
			sb.WriteString(c.Field())
			sb.WriteString(`)`)
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
