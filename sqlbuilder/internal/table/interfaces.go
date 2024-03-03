package table

import (
	"io"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
)

type columnBuilder interface {
	Build() *ast.ColumnSpec
}

type indexBuilder interface {
	Build() *ast.IndexSpec
}

type Formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}
