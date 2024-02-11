package join

import (
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/sel"
)

type JoinBuilder struct {
	kind        ast.JoinKind
	left, right ast.IntoTableExpr
	on          ast.Expr
}

func newJoinBuilder(kind ast.JoinKind, left, right ast.IntoTableExpr) *JoinBuilder {
	return &JoinBuilder{
		kind:  kind,
		left:  left,
		right: right,
	}
}

func Left(left, right string) *JoinBuilder {
	return newJoinBuilder(
		ast.JoinKindLeft,
		sel.Table(left),
		sel.Table(right),
	)
}

func Inner(left, right string) *JoinBuilder {
	return newJoinBuilder(
		ast.JoinKindInner,
		sel.Table(left),
		sel.Table(right),
	)
}

func (b *JoinBuilder) OnEqualColumns(colA, colB string) ast.IntoTableExpr {
	b.on = ast.NewBinaryExpr(
		ast.NewIdentifier(colA),
		ast.BinaryEquals,
		ast.NewIdentifier(colB),
	)

	return ast.NewJoin(
		b.kind,
		b.left.IntoTableExpr(),
		b.right.IntoTableExpr(),
		b.on,
	)
}
