package ast

type TableExpr interface {
	Expr
	tableExpr()
}

type baseTableExpr struct {
	baseExpr
}

type TableName struct {
	baseTableExpr

	Name      string
	Qualifier string
}
