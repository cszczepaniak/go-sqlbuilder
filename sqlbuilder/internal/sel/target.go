package sel

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type Target interface {
	SelectTarget() (string, error)
	ast.IntoTableExpr
}

type Table string

func (t Table) IntoTableExpr() ast.TableExpr {
	return ast.NewTableName(string(t))
}

func (t Table) SelectTarget() (string, error) {
	return string(t), nil
}
