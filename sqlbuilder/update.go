package sqlbuilder

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type updateDialect interface {
	UpdateStmt(table string, fields ...string) (string, error)
	conditioner
}

type fieldAndArg struct {
	field string
	arg   any
}

type UpdateBuilder struct {
	table  string
	fields []fieldAndArg
	upd    updateDialect
	f      filter.Filter
}

func newUpdateBuilder(sel updateDialect, table string) *UpdateBuilder {
	return &UpdateBuilder{
		table: table,
		upd:   sel,
	}
}

func (b *UpdateBuilder) SetFieldTo(field string, val any) *UpdateBuilder {
	b.fields = append(b.fields, fieldAndArg{
		field: field,
		arg:   val,
	})
	return b
}

func (b *UpdateBuilder) Where(f filter.Filter) *UpdateBuilder {
	b.f = f
	return b
}

func (b *UpdateBuilder) WhereAll(f ...filter.Filter) *UpdateBuilder {
	return b.Where(filter.All(f...))
}

func (b *UpdateBuilder) WhereAny(f ...filter.Filter) *UpdateBuilder {
	return b.Where(filter.Any(f...))
}

func (b *UpdateBuilder) Build() (statement.Statement, error) {
	fields, args := b.fieldsAndArgs()

	stmt, err := b.upd.UpdateStmt(b.table, fields...)
	if err != nil {
		return statement.Statement{}, err
	}

	cond, condArgs, err := getCondition(b.upd, b.f)
	if err != nil {
		return statement.Statement{}, err
	}
	stmt += ` ` + cond

	return statement.Statement{
		Stmt: stmt,
		Args: append(args, condArgs...),
	}, nil
}

func (b *UpdateBuilder) Exec(e execer) (sql.Result, error) {
	return exec(b, e)
}

func (b *UpdateBuilder) ExecContext(ctx context.Context, e execCtxer) (sql.Result, error) {
	return execContext(ctx, b, e)
}

func (b *UpdateBuilder) fieldsAndArgs() ([]string, []any) {
	fields := make([]string, 0, len(b.fields))
	args := make([]any, 0, len(b.fields))
	for _, f := range b.fields {
		fields = append(fields, f.field)
		args = append(args, f.arg)
	}
	return fields, args
}
