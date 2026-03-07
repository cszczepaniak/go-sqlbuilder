package sqlbuilder

import (
	"io"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/delete"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/insert"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/sel"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/update"
)

type Formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}

type Builder struct {
	f        Formatter
	database string
}

func New(f Formatter) *Builder {
	return &Builder{
		f: f,
	}
}

func (b *Builder) SetDatabase(db string) *Builder {
	b.database = db
	return b
}

func (b *Builder) qualifiedTableExpr(expr ast.IntoTableExpr) ast.IntoTableExpr {
	if b.database == `` {
		return expr
	}
	qualified := ast.QualifyTableExpr(expr.IntoTableExpr(), b.database)
	return qualified
}

func (b *Builder) SelectFrom(tableExpr ast.IntoTableExpr) *sel.Builder {
	return sel.NewBuilder(b.f, b.qualifiedTableExpr(tableExpr))
}

func (b *Builder) DeleteFrom(tableExpr ast.IntoTableExpr) *delete.Builder {
	return delete.NewBuilder(b.f, b.qualifiedTableExpr(tableExpr))
}

func (b *Builder) Update(tableExpr ast.IntoTableExpr) *update.Builder {
	return update.NewBuilder(b.f, b.qualifiedTableExpr(tableExpr))
}

func (b *Builder) InsertInto(tableExpr ast.IntoTableExpr) *insert.Builder {
	return insert.NewBuilder(b.f, b.qualifiedTableExpr(tableExpr))
}

// CreateTable starts a CREATE TABLE for the given table. It accepts only a bare table
// reference (the result of table.Named("foo")).
func (b *Builder) CreateTable(ref table.BareTableRef) *table.CreateBuilder {
	qualified := b.qualifiedTableExpr(ref)
	name := ast.BaseTableName(qualified.IntoTableExpr())
	return table.NewCreateBuilder(b.f, name)
}
