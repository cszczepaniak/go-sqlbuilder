package sel

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type Target interface {
	ast.IntoTableExpr
}

type Table string

func (t Table) IntoTableExpr() ast.TableExpr {
	return ast.NewTableName(string(t))
}
