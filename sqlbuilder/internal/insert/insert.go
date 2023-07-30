package insert

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/conflict"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/dispatch"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type Dialect interface {
	InsertStmt(table string, fields ...string) (string, error)
	InsertIgnoreStmt(table string, fields ...string) (string, error)
	ValuesStmt(numRecords, argsPerRecord int) (string, error)
	OnConflictStmt(key conflict.Key, conflicts ...conflict.Behavior) (string, error)
}

type Builder struct {
	table     string
	fields    []string
	args      []any
	conflicts *conflictData
	ins       Dialect
}

type conflictData struct {
	key               conflict.Key
	conflictBehaviors []conflict.Behavior
}

func NewBuilder(sel Dialect, table string) *Builder {
	return &Builder{
		table: table,
		ins:   sel,
	}
}

func (b *Builder) Fields(fs ...string) *Builder {
	b.fields = append(b.fields, fs...)
	return b
}

func (b *Builder) Values(vals ...any) *Builder {
	b.args = append(b.args, vals...)
	return b
}

func (b *Builder) OnConflict(key conflict.Key, cs ...conflict.Behavior) *Builder {
	b.conflicts = &conflictData{
		key:               key,
		conflictBehaviors: cs,
	}
	return b
}

func (b *Builder) IgnoreConflicts(key conflict.Key) *Builder {
	c := &conflictData{
		key: key,
	}
	for _, f := range b.fields {
		c.conflictBehaviors = append(c.conflictBehaviors, conflict.Ignore(f))
	}
	b.conflicts = c
	return b
}

func (b *Builder) OverwriteConflicts(key conflict.Key) *Builder {
	c := &conflictData{
		key: key,
	}
	for _, f := range b.fields {
		c.conflictBehaviors = append(c.conflictBehaviors, conflict.Overwrite(f))
	}
	b.conflicts = c
	return b
}

func (b *Builder) Build() (statement.Statement, error) {
	return b.build(b.fields, b.args)
}

func (b *Builder) BuildBatchesOfSize(itemsPerBatch int) ([]statement.Statement, error) {
	if itemsPerBatch <= 0 {
		return nil, errors.New(`batch size must be greater than 0`)
	}
	if err := validate(b.fields, b.args); err != nil {
		return nil, err
	}

	numArgsPerItem := len(b.fields)
	numItems := len(b.args) / numArgsPerItem

	numBatches := (numItems / itemsPerBatch) + 1
	if numItems%itemsPerBatch == 0 {
		numBatches--
	}

	argsPerBatch := itemsPerBatch * numArgsPerItem

	res := make([]statement.Statement, 0, numBatches)
	for i := 0; i < numBatches; i++ {
		start := i * argsPerBatch
		end := start + argsPerBatch

		if end > len(b.args) {
			end = len(b.args)
		}

		stmt, err := b.build(b.fields, b.args[start:end])
		if err != nil {
			return nil, err
		}
		res = append(res, stmt)
	}

	return res, nil
}

func (b *Builder) build(fields []string, args []any) (statement.Statement, error) {
	if err := validate(fields, args); err != nil {
		return statement.Statement{}, err
	}

	stmt, err := b.ins.InsertStmt(b.table, fields...)
	if err != nil {
		return statement.Statement{}, err
	}

	vals, err := b.ins.ValuesStmt(
		len(args)/len(b.fields),
		len(b.fields),
	)
	if err != nil {
		return statement.Statement{}, err
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
			return statement.Statement{}, err
		}
		if conflict != `` {
			stmt += ` ` + conflict
		}
	}

	return statement.Statement{
		Stmt: stmt,
		Args: args,
	}, nil
}

func validate(fields []string, args []any) error {
	if len(fields) == 0 {
		return errors.New(`must provide fields to insert`)
	}
	if len(args)%len(fields) != 0 {
		return errors.New(`number of arguments must be divisible by the number of fields being set`)
	}
	return nil
}

func (b *Builder) Exec(e dispatch.Execer) (sql.Result, error) {
	return dispatch.Exec(b, e)
}

func (b *Builder) ExecContext(ctx context.Context, e dispatch.ExecCtxer) (sql.Result, error) {
	return dispatch.ExecContext(ctx, b, e)
}
