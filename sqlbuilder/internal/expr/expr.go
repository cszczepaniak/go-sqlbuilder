package expr

// Expr represents a generic SQL expression. This can take many forms: a column's name, a boolean expression
// (like Column = ?), a function call (like COUNT(*)), etc. This is a building block for statements. Note that Expr is
// not specific to any SQL dialect; dialects are responsible for synthesizing the specific syntax that corresponds to
// any given expression.
//
// The arguments returned by an expression correspond to placeholder that are included in the resolved expression.
type Expr interface {
	Args() []any
}
