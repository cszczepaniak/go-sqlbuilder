package sqlbuilder

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type builder interface {
	Build() (statement.Statement, error)
}

type queryer interface {
	Query(stmt string, args ...any) (*sql.Rows, error)
}

func query(b builder, q queryer) (*sql.Rows, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return q.Query(res.Stmt, res.Args...)
}

type queryCtxer interface {
	QueryContext(ctx context.Context, stmt string, args ...any) (*sql.Rows, error)
}

func queryContext(ctx context.Context, b builder, q queryCtxer) (*sql.Rows, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return q.QueryContext(ctx, res.Stmt, res.Args...)
}

type rowQueryer interface {
	QueryRow(stmt string, args ...any) *sql.Row
}

func queryRow(b builder, q rowQueryer) (*sql.Row, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return q.QueryRow(res.Stmt, res.Args...), nil
}

type rowQueryCtxer interface {
	QueryRowContext(ctx context.Context, stmt string, args ...any) *sql.Row
}

func queryRowContext(ctx context.Context, b builder, q rowQueryCtxer) (*sql.Row, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return q.QueryRowContext(ctx, res.Stmt, res.Args...), nil
}

type execer interface {
	Exec(stmt string, args ...any) (sql.Result, error)
}

func exec(b builder, e execer) (sql.Result, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return e.Exec(res.Stmt, res.Args...)
}

type execCtxer interface {
	ExecContext(ctx context.Context, stmt string, args ...any) (sql.Result, error)
}

func execContext(ctx context.Context, b builder, e execCtxer) (sql.Result, error) {
	res, err := b.Build()
	if err != nil {
		return nil, err
	}

	return e.ExecContext(ctx, res.Stmt, res.Args...)
}
