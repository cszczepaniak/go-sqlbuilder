package ast

type CreateTable struct {
	Name        *Identifier
	IfNotExists bool

	Columns []*ColumnSpec
}

func NewCreateTable(name string) *CreateTable {
	return &CreateTable{
		Name: NewIdentifier(name),
	}
}

func (c *CreateTable) AcceptVisitor(fn func(n Node) bool) {
	if fn(c) {
		c.Name.AcceptVisitor(fn)
		for _, col := range c.Columns {
			col.AcceptVisitor(fn)
		}
	}
}

func (c *CreateTable) CreateIfNotExists() {
	c.IfNotExists = true
}

func (c *CreateTable) AddColumn(cs *ColumnSpec) {
	c.Columns = append(c.Columns, cs)
}

func (c *CreateTable) PrimaryKey() []*Identifier {
	idents := make([]*Identifier, 0, len(c.Columns))
	for _, col := range c.Columns {
		if col.ComprisesPrimaryKey {
			idents = append(idents, col.Name)
		}
	}
	return idents
}
