package ast

type IntoExpr interface {
	IntoExpr() Expr
}

func IntoExprs(intos ...IntoExpr) []Expr {
	res := make([]Expr, 0, len(intos))
	for _, i := range intos {
		res = append(res, i.IntoExpr())
	}
	return res
}

func None() intoNone {
	return intoNone{}
}

type intoNone struct{}

func (intoNone) IntoExpr() Expr {
	return nil
}

type Expr interface {
	Node
	expr()
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
	BinaryAnd
	BinaryOr
)

type BinaryExpr struct {
	Expr
	Left  Expr
	Op    BinaryExprOperator
	Right Expr
}

func NewBinaryExpr(left IntoExpr, op BinaryExprOperator, right IntoExpr) *BinaryExpr {
	return &BinaryExpr{
		Left:  left.IntoExpr(),
		Op:    op,
		Right: right.IntoExpr(),
	}
}

func (b *BinaryExpr) IntoExpr() Expr {
	return b
}

func (b *BinaryExpr) AcceptVisitor(fn func(Node) bool) {
	if fn(b) {
		b.Left.AcceptVisitor(fn)
		b.Right.AcceptVisitor(fn)
	}
}
