package sel

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/condition"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/dispatch"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/limit"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type Dialect interface {
	SelectStmt(table string, fields ...string) (string, error)
	SelectForUpdateStmt(table string, fields ...string) (string, error)
	OrderBy(o filter.Order) (string, error)

	limit.Limiter
	condition.Conditioner
}

type Builder struct {
	target    Target
	fields    []string
	forUpdate bool
	orderBy   *filter.Order
	sel       Dialect

	*condition.ConditionBuilder[*Builder]
	*limit.LimitBuilder[*Builder]
}

func NewBuilder(sel Dialect, target Target) *Builder {
	b := &Builder{
		target: target,
		sel:    sel,
	}

	b.ConditionBuilder = condition.NewBuilder(b)
	b.LimitBuilder = limit.NewBuilder(b)
	return b
}

func (b *Builder) Fields(fs ...string) *Builder {
	b.fields = append(b.fields, fs...)
	return b
}

func (b *Builder) ForUpdate() *Builder {
	b.forUpdate = true
	return b
}

func (b *Builder) OrderBy(o filter.Order) *Builder {
	b.orderBy = &o
	return b
}

func (b *Builder) Build() (statement.Statement, error) {
	targetStr, err := b.target.SelectTarget()
	if err != nil {
		return statement.Statement{}, err
	}

	var stmt string
	if b.forUpdate {
		stmt, err = b.sel.SelectForUpdateStmt(targetStr, b.fields...)
	} else {
		stmt, err = b.sel.SelectStmt(targetStr, b.fields...)
	}
	if err != nil {
		return statement.Statement{}, err
	}

	cond, args, err := b.ConditionBuilder.SQLAndArgs(b.sel)
	if err != nil {
		return statement.Statement{}, err
	}
	stmt += ` ` + cond

	if b.orderBy != nil {
		order, err := b.sel.OrderBy(*b.orderBy)
		if err != nil {
			return statement.Statement{}, err
		}
		stmt += ` ` + order
	}

	lim, limitArgs, err := b.LimitBuilder.SQLAndArgs(b.sel)
	if err != nil {
		return statement.Statement{}, err
	}
	stmt += ` ` + lim
	args = append(args, limitArgs...)

	return statement.Statement{
		Stmt: stmt,
		Args: args,
	}, nil
}

func (b *Builder) Query(q dispatch.Queryer) (*sql.Rows, error) {
	return dispatch.Query(b, q)
}

func (b *Builder) QueryContext(ctx context.Context, q dispatch.QueryCtxer) (*sql.Rows, error) {
	return dispatch.QueryContext(ctx, b, q)
}

func (b *Builder) QueryRow(q dispatch.RowQueryer) (*sql.Row, error) {
	return dispatch.QueryRow(b, q)
}

func (b *Builder) QueryRowContext(ctx context.Context, q dispatch.RowQueryCtxer) (*sql.Row, error) {
	return dispatch.QueryRowContext(ctx, b, q)
}
