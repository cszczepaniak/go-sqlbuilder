package formatter

import (
	"fmt"
	"io"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
)

type Mysql struct{}

func (m Mysql) FormatNode(w io.Writer, n ast.Node) {
	switch tn := n.(type) {
	case *ast.Select:
		m.formatSelect(w, tn)
	case *ast.Delete:
		m.formatDelete(w, tn)
	case *ast.Insert:
		m.formatInsert(w, tn)
	case *ast.Update:
		m.formatUpdate(w, tn)
	case *ast.TableName:
		m.formatTableName(w, tn)
	case *ast.Identifier:
		m.formatIdentifier(w, tn)
	case *ast.ValuesLiteral:
		m.formatValuesLiteral(w, tn)
	case *ast.Limit:
		m.formatLimit(w, tn)
	case *ast.Lock:
		m.formatLock(w, tn)
	case *ast.Where:
		m.formatWhere(w, tn)
	case *ast.BinaryExpr:
		m.formatBinaryExpr(w, tn)
	case *ast.PlaceholderLiteral:
		m.formatPlaceholderLiteral(w, tn)
	case *ast.TupleLiteral:
		m.formatTupleLiteral(w, tn)
	case *ast.IntegerLiteral:
		m.formatIntegerLiteral(w, tn)
	case *ast.OrderBy:
		m.formatOrderBy(w, tn)
	case *ast.Function:
		m.formatFunction(w, tn)
	case *ast.StarLiteral:
		fmt.Fprint(w, "*")
	case *ast.Distinct:
		m.formatDistinct(w, tn)
	case *ast.OnDuplicateKey:
		m.formatOnDuplicateKey(w, tn)
	default:
		panic(fmt.Sprintf(`unexpected node: %T`, n))
	}
}

func formatCommaDelimited[T ast.Node](w io.Writer, f interface{ FormatNode(w io.Writer, n ast.Node) }, ns ...T) {
	for i, n := range ns {
		f.FormatNode(w, n)
		if i < len(ns)-1 {
			fmt.Fprint(w, `,`)
		}
	}
}

func (m Mysql) formatSelect(w io.Writer, s *ast.Select) {
	fmt.Fprint(w, `SELECT `)
	formatCommaDelimited(w, m, s.Exprs...)

	fmt.Fprint(w, ` FROM `)
	m.FormatNode(w, s.From)

	if s.Where != nil {
		fmt.Fprint(w, ` `)
		m.FormatNode(w, s.Where)
	}
	if s.OrderBy != nil {
		fmt.Fprint(w, ` `)
		m.FormatNode(w, s.OrderBy)
	}
	if s.Limit != nil {
		fmt.Fprint(w, ` `)
		m.FormatNode(w, s.Limit)
	}
	if s.Lock != nil {
		fmt.Fprint(w, ` `)
		m.FormatNode(w, s.Lock)
	}
}

func (m Mysql) formatDelete(w io.Writer, d *ast.Delete) {
	fmt.Fprint(w, `DELETE FROM `)
	m.FormatNode(w, d.From)

	if d.Where != nil {
		fmt.Fprint(w, ` `)
		m.FormatNode(w, d.Where)
	}
	if d.OrderBy != nil {
		fmt.Fprint(w, ` `)
		m.FormatNode(w, d.OrderBy)
	}
	if d.Limit != nil {
		fmt.Fprint(w, ` `)
		m.FormatNode(w, d.Limit)
	}
}

func (m Mysql) formatInsert(w io.Writer, i *ast.Insert) {
	fmt.Fprint(w, `INSERT INTO `)
	m.FormatNode(w, i.Into)
	fmt.Fprint(w, ` (`)
	formatCommaDelimited(w, m, i.Columns...)
	fmt.Fprint(w, `) VALUES `)
	formatCommaDelimited(w, m, i.Values...)
	if i.OnDuplicateKey != nil {
		m.FormatNode(w, i.OnDuplicateKey)
	}
}

func (m Mysql) formatUpdate(w io.Writer, u *ast.Update) {
	fmt.Fprint(w, `UPDATE `)
	m.FormatNode(w, u.Table)
	fmt.Fprint(w, ` SET `)
	formatCommaDelimited(w, m, u.AssignmentList...)

	// TODO we "support" these in the formatter, but we don't expose them to the public via the builders.
	// Add tests for these once we support them publicly.
	if u.OrderBy != nil {
		fmt.Fprint(w, ` ORDER BY `)
		m.FormatNode(w, u.OrderBy)
	}
	if u.Limit != nil {
		fmt.Fprint(w, ` LIMIT `)
		m.FormatNode(w, u.Limit)
	}
}

func (m Mysql) formatFunction(w io.Writer, f *ast.Function) {
	fmt.Fprint(w, f.Name)
	fmt.Fprint(w, `(`)
	formatCommaDelimited(w, m, f.Args...)
	fmt.Fprint(w, `)`)
}

func (m Mysql) formatIntegerLiteral(w io.Writer, l *ast.IntegerLiteral) {
	fmt.Fprintf(w, `%d`, l.Value)
}

func (m Mysql) formatTupleLiteral(w io.Writer, t *ast.TupleLiteral) {
	fmt.Fprint(w, `(`)
	formatCommaDelimited(w, m, t.Values...)
	fmt.Fprint(w, `)`)
}

func (m Mysql) formatPlaceholderLiteral(w io.Writer, _ *ast.PlaceholderLiteral) {
	fmt.Fprint(w, `?`)
}

func (m Mysql) formatWhere(w io.Writer, wh *ast.Where) {
	fmt.Fprint(w, `WHERE `)
	m.FormatNode(w, wh.Expr)
}

func (m Mysql) formatOrderBy(w io.Writer, o *ast.OrderBy) {
	fmt.Fprint(w, `ORDER BY `)
	for i, ord := range o.Orders {
		m.FormatNode(w, ord.Expr)
		switch ord.Direction {
		case ast.OrderAsc:
			fmt.Fprint(w, ` ASC`)
		case ast.OrderDesc:
			fmt.Fprint(w, ` DESC`)
		}

		if i < len(o.Orders)-1 {
			fmt.Fprint(w, `,`)
		}
	}
}

func (m Mysql) formatLock(w io.Writer, l *ast.Lock) {
	if l.Kind == ast.NoLock {
		return
	}
	fmt.Fprint(w, `FOR `)
	switch l.Kind {
	case ast.SharedLock:
		fmt.Fprint(w, `SHARE`)
	case ast.ForUpdateLock:
		fmt.Fprint(w, `UPDATE`)
	}
}

func (m Mysql) formatLimit(w io.Writer, l *ast.Limit) {
	fmt.Fprint(w, `LIMIT `)
	if l.Offset != nil {
		m.FormatNode(w, l.Offset)
		fmt.Fprint(w, `, `)
	}
	m.FormatNode(w, l.Count)
}

func (m Mysql) formatIdentifier(w io.Writer, c *ast.Identifier) {
	fmt.Fprint(w, c.Name)
}

func (m Mysql) formatValuesLiteral(w io.Writer, vl *ast.ValuesLiteral) {
	fmt.Fprint(w, `VALUES(`)
	m.FormatNode(w, vl.Target)
	fmt.Fprint(w, `)`)
}

func (m Mysql) formatTableName(w io.Writer, tn *ast.TableName) {
	if tn.Qualifier != `` {
		fmt.Fprintf(w, `%s.`, tn.Qualifier)
	}
	fmt.Fprint(w, tn.Name)
}

func (m Mysql) formatBinaryExpr(w io.Writer, bin *ast.BinaryExpr) {
	m.FormatNode(w, bin.Left)

	switch bin.Op {
	case ast.BinaryEquals:
		fmt.Fprint(w, ` = `)
	case ast.BinaryNotEquals:
		fmt.Fprint(w, ` != `)
	case ast.BinaryGreater:
		fmt.Fprint(w, ` > `)
	case ast.BinaryGraeaterOrEqual:
		fmt.Fprint(w, ` >= `)
	case ast.BinaryLess:
		fmt.Fprint(w, ` < `)
	case ast.BinaryLessOrEqual:
		fmt.Fprint(w, ` <= `)
	case ast.BinaryIn:
		fmt.Fprint(w, ` IN `)
	}

	m.FormatNode(w, bin.Right)
}

func (m Mysql) formatDistinct(w io.Writer, d *ast.Distinct) {
	fmt.Fprint(w, `DISTINCT `)
	formatCommaDelimited(w, m, d.Exprs...)
}

func (m Mysql) formatOnDuplicateKey(w io.Writer, odk *ast.OnDuplicateKey) {
	fmt.Fprint(w, `ON DUPLICATE KEY UPDATE `)
	formatCommaDelimited(w, m, odk.Updates...)
}
