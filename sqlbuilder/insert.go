package sqlbuilder

import (
	"context"
	"database/sql"
	"errors"
)

type insertDialect interface {
	InsertStmt(table string, fields ...string) (string, error)
	InsertIgnoreStmt(table string, fields ...string) (string, error)
	ValuesStmt(numRecords, argsPerRecord int) (string, error)
}

type InsertBuilder struct {
	table           string
	fields          []string
	args            []any
	ignoreConflicts bool
	ins             insertDialect
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

func (b *InsertBuilder) IgnoreConflicts() *InsertBuilder {
	b.ignoreConflicts = true
	return b
}

func (b *InsertBuilder) Build() (Query, error) {
	if err := b.validate(); err != nil {
		return Query{}, err
	}

	var stmt string
	var err error
	if b.ignoreConflicts {
		stmt, err = b.ins.InsertIgnoreStmt(b.table, b.fields...)
	} else {
		stmt, err = b.ins.InsertStmt(b.table, b.fields...)
	}
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

	return Query{
		Stmt: stmt + ` ` + vals,
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
