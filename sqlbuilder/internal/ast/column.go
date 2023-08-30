package ast

type Column struct {
	Expr
	Name string
}

func NewColumn(name string) *Column {
	return &Column{
		Name: name,
	}
}

func (c *Column) AcceptVisitor(fn func(Node) bool) {
	fn(c)
}

func (c *Column) IntoExpr() Expr {
	return c
}
