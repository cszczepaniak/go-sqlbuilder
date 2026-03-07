package insert

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"strings"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/conflict"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/ast"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/internal/dispatch"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
)

type Formatter interface {
	FormatNode(w io.Writer, n ast.Node)
}

type Builder struct {
	f         Formatter
	table     ast.IntoTableExpr
	columns   []string
	args      []any
	conflicts *conflictData
}

type conflictData struct {
	key               conflict.Key
	conflictBehaviors []conflict.Behavior
}

func NewBuilder(f Formatter, table ast.IntoTableExpr) *Builder {
	return &Builder{
		f:     f,
		table: table,
	}
}

func (b *Builder) Columns(cols ...string) *Builder {
	b.columns = append(b.columns, cols...)
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
	for _, col := range b.columns {
		c.conflictBehaviors = append(c.conflictBehaviors, conflict.Ignore(col))
	}
	b.conflicts = c
	return b
}

func (b *Builder) OverwriteConflicts(key conflict.Key) *Builder {
	c := &conflictData{
		key: key,
	}
	for _, col := range b.columns {
		c.conflictBehaviors = append(c.conflictBehaviors, conflict.Overwrite(col))
	}
	b.conflicts = c
	return b
}

func (b *Builder) Build() (statement.Statement, error) {
	return build(b.f, b.table, b.conflicts, b.columns, b.args)
}

func (b *Builder) BuildBatchesOfSize(itemsPerBatch int) ([]statement.Statement, error) {
	if itemsPerBatch <= 0 {
		return nil, errors.New(`batch size must be greater than 0`)
	}
	if err := validate(b.columns, b.args); err != nil {
		return nil, err
	}

	numArgsPerItem := len(b.columns)
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

		stmt, err := build(b.f, b.table, b.conflicts, b.columns, b.args[start:end])
		if err != nil {
			return nil, err
		}
		res = append(res, stmt)
	}

	return res, nil
}

func build(f Formatter, table ast.IntoTableExpr, conflicts *conflictData, columns []string, args []any) (statement.Statement, error) {
	if err := validate(columns, args); err != nil {
		return statement.Statement{}, err
	}

	idents := make([]*ast.Identifier, 0, len(columns))
	for _, col := range columns {
		idents = append(idents, ast.NewIdentifier(col))
	}

	ins := ast.NewInsert(
		table.IntoTableExpr(),
		idents...,
	)

	for i := 0; i < len(args); i += len(columns) {
		chunk := args[i : i+len(columns)]
		placeholders := make([]ast.IntoExpr, 0, len(chunk))
		for _, arg := range chunk {
			placeholders = append(placeholders, ast.NewPlaceholderLiteral(arg))
		}
		ins.AddValues(placeholders...)
	}

	if conflicts != nil {
		keyIdentNames := conflicts.key.Fields()
		keyIdents := make([]*ast.Identifier, 0, len(keyIdentNames))
		for _, n := range keyIdentNames {
			keyIdents = append(keyIdents, ast.NewIdentifier(n))
		}

		for _, b := range conflicts.conflictBehaviors {
			ins.OnDuplicateKeyUpdate(keyIdents, ast.NewIdentifier(b.Field()), b)
		}
	}

	sb := strings.Builder{}
	f.FormatNode(&sb, ins)

	return statement.Statement{
		Stmt: sb.String(),
		Args: ast.GetArgs(ins),
	}, nil
}

func validate(columns []string, args []any) error {
	if len(columns) == 0 {
		return errors.New(`must provide columns to insert`)
	}
	if len(args)%len(columns) != 0 {
		return errors.New(`number of arguments must be divisible by the number of columns`)
	}
	return nil
}

func (b *Builder) Exec(e dispatch.Execer) (sql.Result, error) {
	return dispatch.Exec(b, e)
}

func (b *Builder) ExecContext(ctx context.Context, e dispatch.ExecCtxer) (sql.Result, error) {
	return dispatch.ExecContext(ctx, b, e)
}
