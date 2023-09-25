package column

import "github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"

type autoIncColumnBuilder[T any, U columnTyper] struct {
	*baseColumnBuilder[T, U]
	autoIncrement bool
}

func newAutoIncColumnBuilder[T any, U columnTyper](name string, parent U) *autoIncColumnBuilder[T, U] {
	return &autoIncColumnBuilder[T, U]{
		baseColumnBuilder: newBaseColumnBuilder[T](name, parent),
	}
}

func (b *autoIncColumnBuilder[T, U]) AutoIncrement() U {
	b.autoIncrement = true
	return b.parent
}

func (b *autoIncColumnBuilder[T, U]) Build() *ast.ColumnSpec {
	cs := b.baseColumnBuilder.Build()
	return cs.SetAutoIncrement(b.autoIncrement)
}

type anyInteger interface {
	int8 | int16 | int32 | int64
}

type integerColumnBuilder[T anyInteger, U columnTyper] struct {
	*autoIncColumnBuilder[T, U]
}

func newIntColumnBuilder[T anyInteger, U columnTyper](name string, parent U) *integerColumnBuilder[T, U] {
	return &integerColumnBuilder[T, U]{
		autoIncColumnBuilder: newAutoIncColumnBuilder[T](name, parent),
	}
}

func (b *integerColumnBuilder[T, U]) Build() *ast.ColumnSpec {
	cs := b.autoIncColumnBuilder.Build()

	if b.defaultVal != nil {
		cs.WithDefault(ast.NewIntegerLiteral(int(*b.defaultVal)))
	}

	return cs
}
