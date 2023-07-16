package sqlbuilder

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/conflict"
)

type insertDialect interface {
	InsertStmt(table string, fields ...string) (string, error)
	InsertIgnoreStmt(table string, fields ...string) (string, error)
	ValuesStmt(numRecords, argsPerRecord int) (string, error)
	OnConflictStmt(conflicts ...conflict.Behavior) (string, error)
}

type InsertBuilder struct {
	table             string
	fields            []string
	args              []any
	conflictBehaviors []conflict.Behavior
	ins               insertDialect
}

func newInsertBuilder(sel insertDialect, table string) *InsertBuilder {
	return &InsertBuilder{
		table: table,
		ins:   sel,
	}
}

func (b *InsertBuilder) Fields(fs ...string) *InsertBuilder {
	b.fields = append(b.fields, fs...)
	return b
}

func (b *InsertBuilder) WithRecord(vals ...any) *InsertBuilder {
	b.args = append(b.args, vals...)
	return b
}

func (b *InsertBuilder) OnConflict(cs ...conflict.Behavior) *InsertBuilder {
	b.conflictBehaviors = append(b.conflictBehaviors, cs...)
	return b
}

func (b *InsertBuilder) IgnoreConflicts() *InsertBuilder {
	for _, f := range b.fields {
		b.conflictBehaviors = append(b.conflictBehaviors, conflict.Ignore(f))
	}
	return b
}

func (b *InsertBuilder) OverwriteConflicts() *InsertBuilder {
	for _, f := range b.fields {
		b.conflictBehaviors = append(b.conflictBehaviors, conflict.Overwrite(f))
	}
	return b
}

func (b *InsertBuilder) Build() (Query, error) {
	if err := b.validate(); err != nil {
		return Query{}, err
	}

	stmt, err := b.ins.InsertStmt(b.table, b.fields...)
	if err != nil {
		return Query{}, err
	}

	vals, err := b.ins.ValuesStmt(
		len(b.args)/len(b.fields),
		len(b.fields),
	)
	if err != nil {
		return Query{}, err
	}
	if vals != `` {
		stmt += ` ` + vals
	}

	conflict, err := b.ins.OnConflictStmt(b.conflictBehaviors...)
	if err != nil {
		return Query{}, err
	}
	if conflict != `` {
		stmt += ` ` + conflict
	}

	return Query{
		Stmt: stmt,
		Args: b.args,
	}, nil
}

func (b *InsertBuilder) Exec(e execer) (sql.Result, error) {
	return exec(b, e)
}

func (b *InsertBuilder) ExecContext(ctx context.Context, e execCtxer) (sql.Result, error) {
	return execContext(ctx, b, e)
}

func (b *InsertBuilder) validate() error {
	if len(b.args)%len(b.fields) != 0 {
		return errors.New(`number of arguments must be divisible by the number of fields being set`)
	}
	return nil
}
