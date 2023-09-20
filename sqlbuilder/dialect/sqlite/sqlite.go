package sqlite

import (
	"fmt"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/functions"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/expr"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/params"
)

func Now() string {
	return `datetime('now')`
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
