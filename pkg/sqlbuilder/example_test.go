package sqlbuilder

import (
	"testing"

	"github.com/cszczepaniak/go-sqlbuilder/pkg/sqlbuilder/dialect/mysql"
	"github.com/cszczepaniak/go-sqlbuilder/pkg/sqlbuilder/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	q, err := New(`Something`).
		Select(mysql.Dialect{}).
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

	q, err = New(`Something`).Delete(mysql.Dialect{}).
		Where(filter.Equals(`A`, 1)).
		Build()
	require.NoError(t, err)

	assert.Len(t, q.Args, 1)
	assert.EqualValues(t, 1, q.Args[0])
	t.Log(q.Stmt)

	q, err = New(`Something`).Update(mysql.Dialect{}).
		SetFieldTo(`A`, 123).
		SetFieldTo(`B`, `foo`).
		Where(filter.Equals(`A`, 1)).
		Build()
	require.NoError(t, err)

	assert.Len(t, q.Args, 3)
	assert.EqualValues(t, []any{123, `foo`, 1}, q.Args)
	t.Log(q.Stmt)

	q, err = New(`Something`).Insert(mysql.Dialect{}).
		Fields(`A`, `B`).
		WithRecord(1, `abc`).
		WithRecord(2, `def`).
		WithRecord(3, `ghi`).
		Build()
	require.NoError(t, err)

	assert.Len(t, q.Args, 6)
	assert.EqualValues(t, []any{1, `abc`, 2, `def`, 3, `ghi`}, q.Args)
	t.Log(q.Stmt)
}
