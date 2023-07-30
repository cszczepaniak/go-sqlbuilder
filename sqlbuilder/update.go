package sqlbuilder

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/condition"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/dispatch"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type updateDialect interface {
	UpdateStmt(table string, fields ...string) (string, error)
	condition.Conditioner
}

type fieldAndArg struct {
	field string
	arg   any
}

type UpdateBuilder struct {
	table  string
	fields []fieldAndArg
	upd    updateDialect

	*condition.ConditionBuilder[*UpdateBuilder]
}

func newUpdateBuilder(sel updateDialect, table string) *UpdateBuilder {
	b := &UpdateBuilder{
		table: table,
		upd:   sel,
	}

	b.ConditionBuilder = condition.NewBuilder(b)
	return b
}

func (b *UpdateBuilder) SetFieldTo(field string, val any) *UpdateBuilder {
	b.fields = append(b.fields, fieldAndArg{
		field: field,
		arg:   val,
	})
	return b
}

func (b *UpdateBuilder) Build() (statement.Statement, error) {
	fields, args := b.fieldsAndArgs()

	stmt, err := b.upd.UpdateStmt(b.table, fields...)
	if err != nil {
		return statement.Statement{}, err
	}

	cond, condArgs, err := b.ConditionBuilder.SQLAndArgs(b.upd)
	if err != nil {
		return statement.Statement{}, err
	}
	stmt += ` ` + cond

	return statement.Statement{
		Stmt: stmt,
		Args: append(args, condArgs...),
	}, nil
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

func (b *UpdateBuilder) Exec(e dispatch.Execer) (sql.Result, error) {
	return dispatch.Exec(b, e)
}

func (b *UpdateBuilder) ExecContext(ctx context.Context, e dispatch.ExecCtxer) (sql.Result, error) {
	return dispatch.ExecContext(ctx, b, e)
}
