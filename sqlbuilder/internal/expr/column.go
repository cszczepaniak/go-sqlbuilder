package expr

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

// Column is a column literal expression. It can optionally be qualified with the database name.
type Column struct {
	Database string
	Name     string
}

func (c Column) IntoExpr() ast.Expr {
	return ast.NewIdentifier(c.Name)
}

func NewColumn(name string) Column {
	return NewQualifiedColumn(``, name)
}

func NewQualifiedColumn(db, name string) Column {
	return Column{
		Database: db,
		Name:     name,
	}
}

func (c Column) IsQualified() bool {
	return c.Database != ``
}

func (c Column) Args() []any {
	return nil
}
