package ast

// IndexSpec specifies an index.
// Example: index_1 (ColA, ColB)
// Example: UNIQUE index_1 (ColA, ColB)
type IndexSpec struct {
	Expr

	Name    *Identifier
	Columns []*Identifier
	Unique  bool
}

func NewIndexSpec(name string, cols ...*Identifier) *IndexSpec {
	return &IndexSpec{
		Name:    NewIdentifier(name),
		Columns: cols,
	}
}

func (is *IndexSpec) AcceptVisitor(fn func(n Node) bool) {
	if fn(is) {
		is.Name.AcceptVisitor(fn)
		for _, col := range is.Columns {
			col.AcceptVisitor(fn)
		}
	}
}
