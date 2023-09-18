package functions

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type Count struct {
	Field    string
	Distinct bool
}

func CountAll() Count {
	return Count{}
}

func CountField(f string) Count {
	return Count{
		Field:    f,
		Distinct: false,
	}
}

func CountDistinct(f string) Count {
	return Count{
		Field:    f,
		Distinct: true,
	}
}

func (c Count) Args() []any {
	return nil
}

func (c Count) All() bool {
	return c.Field == ``
}

func (c Count) IntoExpr() ast.Expr {
	if c.All() {
		return ast.NewFunction(`COUNT`, ast.NewStarLiteral())
	}

	if c.Distinct {
		return ast.NewFunction(`COUNT`, ast.NewDistinct(ast.NewIdentifier(c.Field)))
	}

	return ast.NewFunction(`COUNT`, ast.NewIdentifier(c.Field))
}
