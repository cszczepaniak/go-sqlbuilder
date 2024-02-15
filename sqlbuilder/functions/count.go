package functions

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type Count struct {
	Arg      ast.IntoExpr
	Distinct bool
}

func CountAll() Count {
	return Count{}
}

func CountColumn(name string) Count {
	return Count{
		Arg:      ast.NewIdentifier(name),
		Distinct: false,
	}
}

func CountColumnDistinct(name string) Count {
	return Count{
		Arg:      ast.NewIdentifier(name),
		Distinct: true,
	}
}

func (c Count) Args() []any {
	return nil
}

func (c Count) All() bool {
	return c.Arg == nil
}

func (c Count) IntoExpr() ast.Expr {
	if c.All() {
		return ast.NewFunction(`COUNT`, ast.NewStarLiteral())
	}

	if c.Distinct {
		return ast.NewFunction(`COUNT`, ast.NewDistinct(c.Arg))
	}

	return ast.NewFunction(`COUNT`, c.Arg)
}
