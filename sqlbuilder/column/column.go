package column

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type ColumnExpressionBuilder struct {
	ident *ast.Identifier
	sel   *ast.Selector
	alias *ast.Identifier
}

func (b *ColumnExpressionBuilder) IntoExpr() ast.Expr {
	var expr ast.Expr
	if b.sel != nil {
		expr = b.sel
	} else {
		expr = b.ident
	}

	if b.alias != nil {
		return &ast.Alias{
			ForExpr: expr,
			As:      b.alias,
		}
	}
	return expr
}

func Named(name string) *ColumnExpressionBuilder {
	return &ColumnExpressionBuilder{
		ident: ast.NewIdentifier(name),
	}
}

func (b *ColumnExpressionBuilder) QualifiedBy(qualifier string) *ColumnExpressionBuilder {
	b.sel = &ast.Selector{
		SelectFrom: ast.NewIdentifier(qualifier),
		FieldName:  b.ident,
	}
	return b
}

func (b *ColumnExpressionBuilder) As(alias string) *ColumnExpressionBuilder {
	b.alias = ast.NewIdentifier(alias)
	return b
}
