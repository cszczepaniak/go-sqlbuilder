package ast

type OrderDirection int

const (
	OrderAsc OrderDirection = iota
	OrderDesc
)

type Order struct {
	Expr      Expr
	Direction OrderDirection
}

type OrderBy struct {
	Orders []Order
	Node
}
