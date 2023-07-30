package sqlbuilder

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type selectTarget interface {
	SelectTarget() (string, error)
}

type TableTarget string

func (t TableTarget) SelectTarget() (string, error) {
	return string(t), nil
}

type selectDialect interface {
	SelectStmt(table string, fields ...string) (string, error)
	SelectForUpdateStmt(table string, fields ...string) (string, error)
	OrderBy(o filter.Order) (string, error)
	limiter
	conditioner
}

type SelectBuilder struct {
	target    selectTarget
	fields    []string
	forUpdate bool
	orderBy   *filter.Order
	limit     *int
	sel       selectDialect
	f         filter.Filter
}

func newSelectBuilder(sel selectDialect, target selectTarget) *SelectBuilder {
	return &SelectBuilder{
		target: target,
		sel:    sel,
	}
}

func (b *SelectBuilder) Fields(fs ...string) *SelectBuilder {
	b.fields = append(b.fields, fs...)
	return b
}

func (b *SelectBuilder) ForUpdate() *SelectBuilder {
	b.forUpdate = true
	return b
}

func (b *SelectBuilder) Where(f filter.Filter) *SelectBuilder {
	b.f = f
	return b
}

func (b *SelectBuilder) WhereAll(f ...filter.Filter) *SelectBuilder {
	return b.Where(filter.All(f...))
}

func (b *SelectBuilder) WhereAny(f ...filter.Filter) *SelectBuilder {
	return b.Where(filter.Any(f...))
}

func (b *SelectBuilder) OrderBy(o filter.Order) *SelectBuilder {
	b.orderBy = &o
	return b
}

func (b *SelectBuilder) Limit(limit int) *SelectBuilder {
	b.limit = &limit
	return b
}

func (b *SelectBuilder) Build() (statement.Statement, error) {
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

	cond, args, err := getCondition(b.sel, b.f)
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

	lim, limitArgs, err := getLimit(b.sel, b.limit)
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

func (b *SelectBuilder) Query(q queryer) (*sql.Rows, error) {
	return query(b, q)
}

func (b *SelectBuilder) QueryContext(ctx context.Context, q queryCtxer) (*sql.Rows, error) {
	return queryContext(ctx, b, q)
}

func (b *SelectBuilder) QueryRow(q rowQueryer) (*sql.Row, error) {
	return queryRow(b, q)
}

func (b *SelectBuilder) QueryRowContext(ctx context.Context, q rowQueryCtxer) (*sql.Row, error) {
	return queryRowContext(ctx, b, q)
}
