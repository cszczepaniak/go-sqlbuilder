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
	case *ast.TableName:
		m.formatTableName(w, tn)
	case *ast.Column:
		m.formatColumn(w, tn)
	case *ast.Limit:
		m.formatLimit(w, tn)
	case *ast.Lock:
		m.formatLock(w, tn)
	case ast.BinaryExpr:
		m.formatBinaryExpr(w, tn)
	}
	panic(fmt.Sprintf(`unexpected node: %T`, n))
}

func (m Mysql) formatSelect(w io.Writer, s *ast.Select) {
	fmt.Fprint(w, `SELECT `)
	for i, expr := range s.Exprs {
		m.FormatNode(w, expr)
		if i < len(s.Exprs)-1 {
			fmt.Fprint(w, `,`)
		}
	}

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

func (m Mysql) formatColumn(w io.Writer, c *ast.Column) {
	fmt.Fprint(w, c.Name)
}

func (m Mysql) formatTableName(w io.Writer, tn *ast.TableName) {
	if tn.Qualifier != `` {
		fmt.Fprintf(w, `%s.`, tn.Qualifier)
	}
	fmt.Fprint(w, tn.Name)
}

func (m Mysql) formatBinaryExpr(w io.Writer, bin ast.BinaryExpr) {
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
