package ast

type PlaceholderLiteral struct {
	Expr
	For any
}

func NewPlaceholderLiteral(val any) *PlaceholderLiteral {
	return &PlaceholderLiteral{
		For: val,
	}
}

func (l *PlaceholderLiteral) IntoExpr() Expr {
	return l
}

func (l *PlaceholderLiteral) PlaceholderValue() any {
	return l.For
}

func (l *PlaceholderLiteral) AcceptVisitor(fn func(Node) bool) {
	fn(l)
}

type IntegerLiteral struct {
	Expr
	Value int
}

func NewIntegerLiteral(val int) *IntegerLiteral {
	return &IntegerLiteral{
		Value: val,
	}
}

func (l *IntegerLiteral) IntoExpr() Expr {
	return l
}

func (l *IntegerLiteral) AcceptVisitor(fn func(Node) bool) {
	fn(l)
}

type StarLiteral struct {
	Expr
}

func NewStarLiteral() *StarLiteral {
	return &StarLiteral{}
}

func (l *StarLiteral) IntoExpr() Expr {
	return l
}

func (l *StarLiteral) AcceptVisitor(fn func(Node) bool) {
	fn(l)
}
