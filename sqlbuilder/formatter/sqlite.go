package formatter

import (
	"fmt"
	"io"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
)

type Sqlite struct{}

func (s Sqlite) FormatNode(w io.Writer, n ast.Node) {
	switch tn := n.(type) {
	case *ast.Select:
		s.formatSelect(w, tn)
	case *ast.Delete:
		s.formatDelete(w, tn)
	case *ast.TableName:
		s.formatTableName(w, tn)
	case *ast.Identifier:
		s.formatIdentifier(w, tn)
	case *ast.Limit:
		s.formatLimit(w, tn)
	case *ast.Lock:
		s.formatLock(w, tn)
	case *ast.Where:
		s.formatWhere(w, tn)
	case *ast.BinaryExpr:
		s.formatBinaryExpr(w, tn)
	case *ast.PlaceholderLiteral:
		s.formatPlaceholderLiteral(w, tn)
	case *ast.TupleLiteral:
		s.formatTupleLiteral(w, tn)
	case *ast.IntegerLiteral:
		s.formatIntegerLiteral(w, tn)
	case *ast.OrderBy:
		s.formatOrderBy(w, tn)
	case *ast.Function:
		s.formatFunction(w, tn)
	case *ast.StarLiteral:
		fmt.Fprint(w, "*")
	case *ast.Distinct:
		s.formatDistinct(w, tn)
	default:
		panic(fmt.Sprintf(`unexpected node: %T`, n))
	}
}

func (s Sqlite) formatSelect(w io.Writer, sl *ast.Select) {
	fmt.Fprint(w, `SELECT `)
	for i, expr := range sl.Exprs {
		s.FormatNode(w, expr)
		if i < len(sl.Exprs)-1 {
			fmt.Fprint(w, `,`)
		}
	}

	fmt.Fprint(w, ` FROM `)
	s.FormatNode(w, sl.From)

	if sl.Where != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, sl.Where)
	}
	if sl.OrderBy != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, sl.OrderBy)
	}
	if sl.Limit != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, sl.Limit)
	}
	if sl.Lock != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, sl.Lock)
	}
}

func (s Sqlite) formatDelete(w io.Writer, d *ast.Delete) {
	fmt.Fprint(w, `DELETE FROM `)
	s.FormatNode(w, d.From)

	if d.Where != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, d.Where)
	}
	if d.OrderBy != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, d.OrderBy)
	}
	if d.Limit != nil {
		fmt.Fprint(w, ` `)
		s.FormatNode(w, d.Limit)
	}
}

func (s Sqlite) formatFunction(w io.Writer, f *ast.Function) {
	fmt.Fprint(w, f.Name)
	fmt.Fprint(w, `(`)
	for _, arg := range f.Args {
		s.FormatNode(w, arg)
	}
	fmt.Fprint(w, `)`)
}

func (s Sqlite) formatIntegerLiteral(w io.Writer, l *ast.IntegerLiteral) {
	fmt.Fprintf(w, `%d`, l.Value)
}

func (s Sqlite) formatTupleLiteral(w io.Writer, t *ast.TupleLiteral) {
	fmt.Fprint(w, `(`)
	for i, val := range t.Values {
		s.FormatNode(w, val)
		if i < len(t.Values)-1 {
			fmt.Fprint(w, `,`)
		}
	}
	fmt.Fprint(w, `)`)
}

func (s Sqlite) formatPlaceholderLiteral(w io.Writer, _ *ast.PlaceholderLiteral) {
	fmt.Fprint(w, `?`)
}

func (s Sqlite) formatWhere(w io.Writer, wh *ast.Where) {
	fmt.Fprint(w, `WHERE `)
	s.FormatNode(w, wh.Expr)
}

func (s Sqlite) formatOrderBy(w io.Writer, o *ast.OrderBy) {
	fmt.Fprint(w, `ORDER BY `)
	for i, ord := range o.Orders {
		s.FormatNode(w, ord.Expr)
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

func (s Sqlite) formatLock(w io.Writer, l *ast.Lock) {
	// SQLite doesn't support locking. It doesn't support multiple concurrent transactions, so every select is more or less equivalent to MySQL's FOR UPDATE.
}

func (s Sqlite) formatLimit(w io.Writer, l *ast.Limit) {
	fmt.Fprint(w, `LIMIT `)
	if l.Offset != nil {
		s.FormatNode(w, l.Offset)
		fmt.Fprint(w, `, `)
	}
	s.FormatNode(w, l.Count)
}

func (s Sqlite) formatIdentifier(w io.Writer, c *ast.Identifier) {
	fmt.Fprint(w, c.Name)
}

func (s Sqlite) formatTableName(w io.Writer, tn *ast.TableName) {
	if tn.Qualifier != `` {
		fmt.Fprintf(w, `%s.`, tn.Qualifier)
	}
	fmt.Fprint(w, tn.Name)
}

func (s Sqlite) formatBinaryExpr(w io.Writer, bin *ast.BinaryExpr) {
	s.FormatNode(w, bin.Left)

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

	s.FormatNode(w, bin.Right)
}

func (s Sqlite) formatDistinct(w io.Writer, d *ast.Distinct) {
	fmt.Fprint(w, `DISTINCT `)
	formatCommaDelimited(w, s, d.Exprs...)
}
