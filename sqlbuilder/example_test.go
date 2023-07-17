package sqlbuilder

import (
	"testing"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/dialect/mysql"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	q, err := New(mysql.Dialect{}).
		Select("Something").
		Fields(`A`, `B`).
		WhereAll(
			filter.In(`A`, 1, 2, 3),
			filter.Any(
				filter.Equals(`B`, "abc"),
				filter.Equals(`B`, "def"),
			),
		).
		Build()
	require.NoError(t, err)

	assert.Len(t, q.Args, 5)
	assert.EqualValues(t, []any{1, 2, 3, "abc", "def"}, q.Args)
	t.Log(q.Stmt)

	q, err = New(mysql.Dialect{}).Delete("Something").
		Where(filter.Equals(`A`, 1)).
		Build()
	require.NoError(t, err)

	assert.Len(t, q.Args, 1)
	assert.EqualValues(t, 1, q.Args[0])
	t.Log(q.Stmt)

	q, err = New(mysql.Dialect{}).Update("Something").
		SetFieldTo(`A`, 123).
		SetFieldTo(`B`, `foo`).
		Where(filter.Equals(`A`, 1)).
		Build()
	require.NoError(t, err)

	assert.Len(t, q.Args, 3)
	assert.EqualValues(t, []any{123, `foo`, 1}, q.Args)
	t.Log(q.Stmt)

	q, err = New(mysql.Dialect{}).Insert("Something").
		Fields(`A`, `B`).
		Values(1, `abc`).
		Values(2, `def`).
		Values(3, `ghi`).
		Build()
	require.NoError(t, err)

	assert.Len(t, q.Args, 6)
	assert.EqualValues(t, []any{1, `abc`, 2, `def`, 3, `ghi`}, q.Args)
	t.Log(q.Stmt)
}
