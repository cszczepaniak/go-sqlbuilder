package delete

import (
	"context"
	"database/sql"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/condition"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/dispatch"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/limit"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type Dialect interface {
	DeleteStmt(table string) (string, error)

	limit.Limiter
	condition.Conditioner
}

type Builder struct {
	table string
	del   Dialect

	*condition.ConditionBuilder[*Builder]
	*limit.LimitBuilder[*Builder]
}

func NewBuilder(sel Dialect, table string) *Builder {
	b := &Builder{
		table: table,
		del:   sel,
	}
	b.ConditionBuilder = condition.NewBuilder(b)
	b.LimitBuilder = limit.NewBuilder(b)
	return b
}

func (b *Builder) Build() (statement.Statement, error) {
	stmt, err := b.del.DeleteStmt(b.table)
	if err != nil {
		return statement.Statement{}, err
	}

	cond, args, err := b.ConditionBuilder.SQLAndArgs(b.del)
	if err != nil {
		return statement.Statement{}, err
	}
	stmt += ` ` + cond

	lim, limitArgs, err := b.LimitBuilder.SQLAndArgs(b.del)
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

func (b *Builder) Exec(e dispatch.Execer) (sql.Result, error) {
	return dispatch.Exec(b, e)
}

func (b *Builder) ExecContext(ctx context.Context, e dispatch.ExecCtxer) (sql.Result, error) {
	return dispatch.ExecContext(ctx, b, e)
}
