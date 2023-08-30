package ast

type IntoExpr interface {
	IntoExpr() Expr
}

type Expr interface {
	Node
	expr()
}

type baseExpr struct {
	Expr
}

func (be baseExpr) IntoExpr() Expr {
	return be
}

type BinaryExprOperator int

const (
	BinaryEquals BinaryExprOperator = iota
	BinaryNotEquals
	BinaryGreater
	BinaryGraeaterOrEqual
	BinaryLess
	BinaryLessOrEqual
	BinaryIn
)

type BinaryExpr struct {
	Expr
	Left  Expr
	Op    BinaryExprOperator
	Right Expr
}
