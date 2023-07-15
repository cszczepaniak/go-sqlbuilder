package sqlbuilder

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/pkg/sqlbuilder/filter"
)

type deleteDialect interface {
	DeleteStmt(table string) (string, error)
	conditioner
}

type DeleteBuilder struct {
	table string
	del   deleteDialect
	f     filter.Filter
}

func newDeleteBuilder(sel deleteDialect, table string) *DeleteBuilder {
	return &DeleteBuilder{
		table: table,
		del:   sel,
	}
}

func (b *DeleteBuilder) Where(f filter.Filter) *DeleteBuilder {
	b.f = f
	return b
}

func (b *DeleteBuilder) WhereAll(f ...filter.Filter) *DeleteBuilder {
	return b.Where(filter.All(f...))
}

func (b *DeleteBuilder) WhereAny(f ...filter.Filter) *DeleteBuilder {
	return b.Where(filter.Any(f...))
}

func (b *DeleteBuilder) Build() (Query, error) {
	stmt, err := b.del.DeleteStmt(b.table)
	if err != nil {
		return Query{}, err
	}

	cond, args, err := getCondition(b.del, b.f)
	if err != nil {
		return Query{}, err
	}
	stmt += ` ` + cond

	return Query{
		Stmt: stmt,
		Args: args,
	}, nil
}

func (b *DeleteBuilder) Exec(e execer) (sql.Result, error) {
	return exec(b, e)
}

func (b *DeleteBuilder) ExecContext(ctx context.Context, e execCtxer) (sql.Result, error) {
	return execContext(ctx, b, e)
}
