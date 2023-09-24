package ast

type ColumnSpec struct {
	Name             *Identifier
	Type             ColumnType
	Nullability      Nullability
	Default          *ColumnDefault
	AutoIncrementing *AutoIncrement

	ComprisesPrimaryKey bool
}

func NewColumnSpec(name string, typ ColumnType) *ColumnSpec {
	return &ColumnSpec{
		Name:        NewIdentifier(name),
		Type:        typ,
		Nullability: NoNullability,
	}
}

func (cs *ColumnSpec) AcceptVisitor(fn func(n Node) bool) {
	if fn(cs) {
		cs.Name.AcceptVisitor(fn)
		cs.Type.AcceptVisitor(fn)
		if cs.Nullability != NoNullability {
			cs.Nullability.AcceptVisitor(fn)
		}
		if cs.Default != nil {
			cs.Default.AcceptVisitor(fn)
		}
		if cs.AutoIncrementing != nil {
			cs.AutoIncrementing.AcceptVisitor(fn)
		}
	}
}

func (c *ColumnSpec) NotNull() {
	c.Nullability = NotNull
}

func (c *ColumnSpec) Null() {
	c.Nullability = Null
}

func (c *ColumnSpec) AutoIncrement() {
	c.AutoIncrementing = &AutoIncrement{}
}

func (c *ColumnSpec) WithDefault(val any) {
	c.Default = newColumnDefault(val)
}

func (c *ColumnSpec) PrimaryKey() {
	c.ComprisesPrimaryKey = true
}
