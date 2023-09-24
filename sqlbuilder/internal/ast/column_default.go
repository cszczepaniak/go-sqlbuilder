package ast

type ColumnDefault struct {
	Value any
}

func newColumnDefault(val any) *ColumnDefault {
	return &ColumnDefault{
		Value: val,
	}
}

func (c *ColumnDefault) AcceptVisitor(fn func(n Node) bool) {
	fn(c)
}
