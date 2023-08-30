package filter

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type Filter interface {
	Args() []any
	ast.IntoExpr
}

type AllFilter struct {
	Filters []Filter
}

func All(fs ...Filter) AllFilter {
	return AllFilter{
		Filters: fs,
	}
}

func (f AllFilter) Args() []any {
	var args []any
	for _, ff := range f.Filters {
		args = append(args, ff.Args()...)
	}
	return args
}

func (f AllFilter) IntoExpr() ast.Expr {
	return makeChainedExpr(f.Filters[0], ast.BinaryAnd, f.Filters[1:]...).IntoExpr()
}

func makeChainedExpr(left Filter, op ast.BinaryExprOperator, rest ...Filter) ast.IntoExpr {
	if len(rest) == 0 {
		return left
	}
	return ast.NewBinaryExpr(left, op, makeChainedExpr(rest[0], op, rest[1:]...))
}

type AnyFilter struct {
	Filters []Filter
}

func Any(fs ...Filter) AnyFilter {
	return AnyFilter{
		Filters: fs,
	}
}

func (f AnyFilter) Args() []any {
	var args []any
	for _, ff := range f.Filters {
		args = append(args, ff.Args()...)
	}
	return args
}

func (f AnyFilter) IntoExpr() ast.Expr {
	return makeChainedExpr(f.Filters[0], ast.BinaryOr, f.Filters[1:]...).IntoExpr()
}

type EqualsFilter struct {
	Column string
	Value  any
}

func Equals(column string, val any) EqualsFilter {
	return EqualsFilter{
		Column: column,
		Value:  val,
	}
}

func (f EqualsFilter) Args() []any {
	return []any{f.Value}
}

func (f EqualsFilter) IntoExpr() ast.Expr {
	return ast.NewBinaryExpr(ast.NewColumn(f.Column), ast.BinaryEquals, ast.NewPlaceholderLiteral(f.Value))
}

type NotEqualsFilter struct {
	Column string
	Value  any
}

func NotEquals(column string, val any) NotEqualsFilter {
	return NotEqualsFilter{
		Column: column,
		Value:  val,
	}
}

func (f NotEqualsFilter) Args() []any {
	return []any{f.Value}
}

func (f NotEqualsFilter) IntoExpr() ast.Expr {
	return ast.NewBinaryExpr(ast.NewColumn(f.Column), ast.BinaryNotEquals, ast.NewPlaceholderLiteral(f.Value))
}

type GreaterFilter struct {
	Column string
	Value  any
}

func Greater(column string, val any) GreaterFilter {
	return GreaterFilter{
		Column: column,
		Value:  val,
	}
}

func (f GreaterFilter) Args() []any {
	return []any{f.Value}
}

func (f GreaterFilter) IntoExpr() ast.Expr {
	return ast.NewBinaryExpr(ast.NewColumn(f.Column), ast.BinaryGreater, ast.NewPlaceholderLiteral(f.Value))
}

type GreaterOrEqualFilter struct {
	Column string
	Value  any
}

func GreaterOrEqual(column string, val any) GreaterOrEqualFilter {
	return GreaterOrEqualFilter{
		Column: column,
		Value:  val,
	}
}

func (f GreaterOrEqualFilter) Args() []any {
	return []any{f.Value}
}

func (f GreaterOrEqualFilter) IntoExpr() ast.Expr {
	return ast.NewBinaryExpr(ast.NewColumn(f.Column), ast.BinaryGraeaterOrEqual, ast.NewPlaceholderLiteral(f.Value))
}

type LessFilter struct {
	Column string
	Value  any
}

func Less(column string, val any) LessFilter {
	return LessFilter{
		Column: column,
		Value:  val,
	}
}

func (f LessFilter) Args() []any {
	return []any{f.Value}
}

func (f LessFilter) IntoExpr() ast.Expr {
	return ast.NewBinaryExpr(ast.NewColumn(f.Column), ast.BinaryLess, ast.NewPlaceholderLiteral(f.Value))
}

type LessOrEqualFilter struct {
	Column string
	Value  any
}

func LessOrEqual(column string, val any) LessOrEqualFilter {
	return LessOrEqualFilter{
		Column: column,
		Value:  val,
	}
}

func (f LessOrEqualFilter) Args() []any {
	return []any{f.Value}
}

func (f LessOrEqualFilter) IntoExpr() ast.Expr {
	return ast.NewBinaryExpr(ast.NewColumn(f.Column), ast.BinaryLessOrEqual, ast.NewPlaceholderLiteral(f.Value))
}

type InFilter struct {
	Column string
	Values []any
}

func In(column string, vals ...any) InFilter {
	return InFilter{
		Column: column,
		Values: vals,
	}
}

func (f InFilter) Args() []any {
	return f.Values
}

func (f InFilter) IntoExpr() ast.Expr {
	exprs := make([]ast.IntoExpr, 0, len(f.Values))
	for _, val := range f.Values {
		exprs = append(exprs, ast.NewPlaceholderLiteral(val))
	}
	return ast.NewBinaryExpr(ast.NewColumn(f.Column), ast.BinaryIn, ast.NewTupleLiteral(exprs...))
}
