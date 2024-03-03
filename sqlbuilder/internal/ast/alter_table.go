package ast

type AlterTable struct {
	Name *Identifier

	AddColumns []*ColumnSpec
	AddIndices []*IndexSpec
}

func NewAlterTable(name string) *AlterTable {
	return &AlterTable{
		Name: NewIdentifier(name),
	}
}

func (at *AlterTable) AcceptVisitor(fn func(n Node) bool) {
	if fn(at) {
		at.Name.AcceptVisitor(fn)
		for _, col := range at.AddColumns {
			col.AcceptVisitor(fn)
		}
		for _, idx := range at.AddIndices {
			idx.AcceptVisitor(fn)
		}
	}
}

func (at *AlterTable) AddColumn(cs *ColumnSpec) {
	at.AddColumns = append(at.AddColumns, cs)
}

func (at *AlterTable) AddIndex(is *IndexSpec) {
	at.AddIndices = append(at.AddIndices, is)
}
