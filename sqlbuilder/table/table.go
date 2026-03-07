package table

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

// BareTableRef is implemented only by the result of Named(name). It is accepted by
// CreateTable so that you can pass the same table reference, but the type system
// rejects joins and aliased tables (e.g. Named("x").As("y")) at compile time.
type BareTableRef interface {
	ast.IntoTableExpr
	bareTableRef()
}

// BareTable represents a single table name with no alias or join. Returned by Named(name).
type BareTable struct {
	name string
}

func (b *BareTable) IntoTableExpr() ast.TableExpr {
	return ast.NewTableName(b.name)
}

func (b *BareTable) bareTableRef() {}

// Named returns a bare table reference for the given table name. Use it with
// SelectFrom, InsertInto, Update, DeleteFrom, and CreateTable. Call As(alias) to
// get a general table expr for joins.
func Named(name string) *BareTable {
	return &BareTable{name: name}
}

// As returns a TableBuilder with an alias, for use in joins. The result does not
// implement BareTableRef, so it cannot be passed to CreateTable.
func (b *BareTable) As(alias string) *TableBuilder {
	return newTableBuilder(&ast.TableAlias{
		ForExpr: ast.NewTableName(b.name),
		As:      ast.NewIdentifier(alias),
	})
}

// LeftJoin starts a left join from this table.
func (b *BareTable) LeftJoin(tableExpr ast.IntoTableExpr) *JoinBuilder {
	return &JoinBuilder{
		kind:       ast.JoinKindLeft,
		joiningTo:  b,
		toBeJoined: tableExpr,
	}
}

// InnerJoin starts an inner join from this table.
func (b *BareTable) InnerJoin(tableExpr ast.IntoTableExpr) *JoinBuilder {
	return &JoinBuilder{
		kind:       ast.JoinKindInner,
		joiningTo:  b,
		toBeJoined: tableExpr,
	}
}

type TableBuilder struct {
	tableExpr ast.IntoTableExpr
}

func newTableBuilder(tableExpr ast.IntoTableExpr) *TableBuilder {
	return &TableBuilder{
		tableExpr: tableExpr,
	}
}

func (tb *TableBuilder) IntoTableExpr() ast.TableExpr {
	return tb.tableExpr.IntoTableExpr()
}

func (tb *TableBuilder) As(alias string) *TableBuilder {
	tb.tableExpr = &ast.TableAlias{
		ForExpr: tb.tableExpr.IntoTableExpr(),
		As:      ast.NewIdentifier(alias),
	}
	return tb
}

func (tb *TableBuilder) LeftJoin(tableExpr ast.IntoTableExpr) *JoinBuilder {
	return &JoinBuilder{
		kind:       ast.JoinKindLeft,
		joiningTo:  tb,
		toBeJoined: tableExpr,
	}
}

func (tb *TableBuilder) InnerJoin(tableExpr ast.IntoTableExpr) *JoinBuilder {
	return &JoinBuilder{
		kind:       ast.JoinKindInner,
		joiningTo:  tb,
		toBeJoined: tableExpr,
	}
}
