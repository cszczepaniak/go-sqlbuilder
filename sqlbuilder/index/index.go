package index

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type IndexBuilder struct {
	idx *ast.IndexSpec
}

func New(name string) *IndexBuilder {
	return &IndexBuilder{
		idx: ast.NewIndexSpec(name),
	}
}

type columnNamer interface {
	ColumnName() string
}

func (b *IndexBuilder) OnColumns(cs ...columnNamer) *IndexBuilder {
	for _, c := range cs {
		b.idx.Columns = append(b.idx.Columns, ast.NewIdentifier(c.ColumnName()))
	}
	return b
}

func (b *IndexBuilder) Build() *ast.IndexSpec {
	return b.idx
}
