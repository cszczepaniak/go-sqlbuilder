package ast

type Limit struct {
	Node
	Offset Expr
	Count  Expr
}
