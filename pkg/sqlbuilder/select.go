package sqlbuilder

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/pkg/sqlbuilder/filter"
)

type selectDialect interface {
	SelectStmt(table string, fields ...string) (string, error)
	SelectForUpdateStmt(table string, fields ...string) (string, error)
	OrderBy(o filter.Order) (string, error)
	conditioner
}

type SelectBuilder struct {
	table     string
	fields    []string
	forUpdate bool
	orderBy   *filter.Order
	sel       selectDialect
	f         filter.Filter
}

func newSelectBuilder(sel selectDialect, table string) *SelectBuilder {
	return &SelectBuilder{
		table: table,
		sel:   sel,
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

func (b *SelectBuilder) Build() (Query, error) {
	var stmt string
	var err error
	if b.forUpdate {
		stmt, err = b.sel.SelectForUpdateStmt(b.table, b.fields...)
	} else {
		stmt, err = b.sel.SelectStmt(b.table, b.fields...)
	}
	if err != nil {
		return Query{}, err
	}

	cond, args, err := getCondition(b.sel, b.f)
	if err != nil {
		return Query{}, err
	}
	stmt += ` ` + cond

	if b.orderBy != nil {
		order, err := b.sel.OrderBy(*b.orderBy)
		if err != nil {
			return Query{}, err
		}
		stmt += ` ` + order
	}

	return Query{
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
