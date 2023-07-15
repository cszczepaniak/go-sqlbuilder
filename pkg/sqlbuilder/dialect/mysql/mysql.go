package mysql

import (
	"fmt"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/pkg/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/pkg/sqlbuilder/internal/params"
)

func Now() string {
	return `NOW()`
}

type Dialect struct{}

func (m Dialect) SelectStmt(table string, fields ...string) (string, error) {
	return `SELECT ` + strings.Join(fields, `,`) + ` FROM ` + table, nil
}

func (m Dialect) SelectForUpdateStmt(table string, fields ...string) (string, error) {
	return `SELECT ` + strings.Join(fields, `,`) + ` FROM ` + table + ` FOR UPDATE`, nil
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

func (m Dialect) Terminator() string {
	return `;`
}
