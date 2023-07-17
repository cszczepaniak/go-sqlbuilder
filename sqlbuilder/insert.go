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
	OnConflictStmt(key conflict.Key, conflicts ...conflict.Behavior) (string, error)
}

type InsertBuilder struct {
	table     string
	fields    []string
	args      []any
	conflicts *conflictData
	ins       insertDialect
}

type conflictData struct {
	key               conflict.Key
	conflictBehaviors []conflict.Behavior
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

func (b *InsertBuilder) OnConflict(key conflict.Key, cs ...conflict.Behavior) *InsertBuilder {
	b.conflicts = &conflictData{
		key:               key,
		conflictBehaviors: cs,
	}
	return b
}

func (b *InsertBuilder) IgnoreConflicts(key conflict.Key) *InsertBuilder {
	c := &conflictData{
		key: key,
	}
	for _, f := range b.fields {
		c.conflictBehaviors = append(c.conflictBehaviors, conflict.Ignore(f))
	}
	b.conflicts = c
	return b
}

func (b *InsertBuilder) OverwriteConflicts(key conflict.Key) *InsertBuilder {
	c := &conflictData{
		key: key,
	}
	for _, f := range b.fields {
		c.conflictBehaviors = append(c.conflictBehaviors, conflict.Overwrite(f))
	}
	b.conflicts = c
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

	if b.conflicts != nil {
		conflict, err := b.ins.OnConflictStmt(
			b.conflicts.key,
			b.conflicts.conflictBehaviors...,
		)
		if err != nil {
			return Query{}, err
		}
		if conflict != `` {
			stmt += ` ` + conflict
		}
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
