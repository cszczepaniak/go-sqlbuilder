package ast

type IntoTableExpr interface {
	IntoTableExpr() TableExpr
}

type TableExpr interface {
	Expr
	tableExpr()
}

type TableName struct {
	TableExpr
	*Identifier
}

func NewTableName(name string) *TableName {
	return &TableName{
		Identifier: NewIdentifier(name),
	}
}

func (t *TableName) IntoTableExpr() TableExpr {
	return t
}

func (t *TableName) AcceptVisitor(fn func(Node) bool) {
	fn(t)
}

// BaseTableName returns the underlying table name for a simple table reference.
// For TableName it returns the name; for TableAlias it recurses on the inner expr.
// For Join it returns the base name of the left side.
func BaseTableName(expr TableExpr) string {
	switch e := expr.(type) {
	case *TableName:
		return e.Name
	case *TableAlias:
		return BaseTableName(e.ForExpr)
	case *Join:
		return BaseTableName(e.Left)
	default:
		return ""
	}
}

// QualifyTableExpr returns a new TableExpr with the leading table name prefixed by qualifier.
// e.g. QualifyTableExpr(TableName("users"), "mydb") -> TableName("mydb.users").
func QualifyTableExpr(expr TableExpr, qualifier string) TableExpr {
	if qualifier == "" {
		return expr
	}
	switch e := expr.(type) {
	case *TableName:
		return NewTableName(qualifier + "." + e.Name)
	case *TableAlias:
		return &TableAlias{
			ForExpr: QualifyTableExpr(e.ForExpr, qualifier),
			As:      e.As,
		}
	case *Join:
		return NewJoin(e.Kind, QualifyTableExpr(e.Left, qualifier), e.Right, e.On)
	default:
		return expr
	}
}

type TupleLiteral struct {
	TableExpr

	Values []Expr
}

func NewTupleLiteral(values ...IntoExpr) *TupleLiteral {
	return &TupleLiteral{
		Values: IntoExprs(values...),
	}
}

func (t *TupleLiteral) IntoExpr() Expr {
	return t
}

func (t *TupleLiteral) AcceptVisitor(fn func(Node) bool) {
	if fn(t) {
		for _, v := range t.Values {
			v.AcceptVisitor(fn)
		}
	}
}
