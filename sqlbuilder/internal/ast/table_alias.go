package ast

// TableAlias represents "table_expr AS alias". ForExpr is always a TableExpr
// (the thing being aliased), so no type assertion is needed when recursing.
type TableAlias struct {
	TableExpr
	ForExpr TableExpr
	As      *Identifier
}

func (a *TableAlias) IntoTableExpr() TableExpr {
	return a
}

func (a *TableAlias) AcceptVisitor(fn func(Node) bool) {
	fn(a)
}
