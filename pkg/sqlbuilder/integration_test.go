package sqlbuilder_test

import (
	"database/sql"
	"os"
	"path"
	"testing"

	"github.com/cszczepaniak/go-sqlbuilder/pkg/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/pkg/sqlbuilder/dialect/sqlite"
	"github.com/cszczepaniak/go-sqlbuilder/pkg/sqlbuilder/filter"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openSQLiteDatabase(t *testing.T) *sql.DB {
	t.Helper()

	dir, err := os.MkdirTemp(``, ``)
	require.NoError(t, err)

	dataSource := path.Join(dir, `sqlite-database.db`)

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll(dir))
	})

	db, err := sql.Open(`sqlite3`, dataSource)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, db.Close())
	})

	return db
}

func TestSQLite(t *testing.T) {
	db := openSQLiteDatabase(t)

	_, err := db.Exec(`CREATE TABLE Example (
		ID TEXT NOT NULL PRIMARY KEY,
		NumberField INT,
		TextField TEXT
	)`)
	require.NoError(t, err)

	b := sqlbuilder.New(`Example`)

	res, err := b.Insert(sqlite.Dialect{}).
		Fields(`ID`, `NumberField`, `TextField`).
		WithRecord(`a`, 1, `aa`).
		WithRecord(`b`, 2, `bb`).
		WithRecord(`c`, 3, `cc`).
		WithRecord(`d`, 4, `dd`).
		WithRecord(`e`, 5, `ee`).
		Exec(db)
	require.NoError(t, err)

	n, err := res.RowsAffected()
	require.NoError(t, err)
	assert.EqualValues(t, 5, n)

	row, err := b.Select(sqlite.Dialect{}).
		Fields(`NumberField`, `TextField`).
		Where(filter.Equals(`NumberField`, 3)).
		QueryRow(db)
	require.NoError(t, err)

	{
		var (
			numberField int
			textField   string
		)
		err = row.Scan(&numberField, &textField)
		require.NoError(t, err)
		assert.Equal(t, 3, numberField)
		assert.Equal(t, `cc`, textField)
	}

	rows, err := b.Select(sqlite.Dialect{}).
		Fields(`ID`, `NumberField`, `TextField`).
		Where(filter.In(`TextField`, `bb`, `dd`)).
		Query(db)
	require.NoError(t, err)

	{
		var (
			id          string
			numberField int
			textField   string
		)

		assert.True(t, rows.Next())
		require.NoError(t, rows.Scan(&id, &numberField, &textField))
		assert.Equal(t, `b`, id)
		assert.Equal(t, 2, numberField)
		assert.Equal(t, `bb`, textField)

		assert.True(t, rows.Next())
		require.NoError(t, rows.Scan(&id, &numberField, &textField))
		assert.Equal(t, `d`, id)
		assert.Equal(t, 4, numberField)
		assert.Equal(t, `dd`, textField)

		assert.False(t, rows.Next())
	}

	res, err = b.Update(sqlite.Dialect{}).
		SetFieldTo(`NumberField`, 123).
		SetFieldTo(`TextField`, `gotcha`).
		WhereAll(
			filter.NotEquals(`TextField`, `bb`),
			filter.LessOrEqual(`NumberField`, 2),
		).
		Exec(db)
	require.NoError(t, err)

	n, err = res.RowsAffected()
	require.NoError(t, err)
	assert.EqualValues(t, 1, n)

	row, err = b.Select(sqlite.Dialect{}).
		Fields(`*`).
		Where(filter.Equals(`ID`, `a`)).
		QueryRow(db)
	require.NoError(t, err)

	{
		var (
			id          string
			numberField int
			textField   string
		)
		err = row.Scan(&id, &numberField, &textField)
		require.NoError(t, err)
		assert.Equal(t, `a`, id)
		assert.Equal(t, 123, numberField)
		assert.Equal(t, `gotcha`, textField)
	}

	// It works with transactions too.
	tx, err := db.Begin()
	require.NoError(t, err)

	res, err = b.Delete(sqlite.Dialect{}).
		Where(filter.Greater(`NumberField`, 3)).
		Exec(tx)
	require.NoError(t, err)

	n, err = res.RowsAffected()
	require.NoError(t, err)
	assert.EqualValues(t, 3, n)

	require.NoError(t, tx.Commit())

	rows, err = b.Select(sqlite.Dialect{}).
		Fields(`ID`, `NumberField`, `TextField`).
		OrderBy(filter.OrderDesc(`NumberField`)).
		Query(db)
	require.NoError(t, err)

	{
		var (
			id          string
			numberField int
			textField   string
		)

		assert.True(t, rows.Next())
		require.NoError(t, rows.Scan(&id, &numberField, &textField))
		assert.Equal(t, `c`, id)
		assert.Equal(t, 3, numberField)
		assert.Equal(t, `cc`, textField)

		assert.True(t, rows.Next())
		require.NoError(t, rows.Scan(&id, &numberField, &textField))
		assert.Equal(t, `b`, id)
		assert.Equal(t, 2, numberField)
		assert.Equal(t, `bb`, textField)

		assert.False(t, rows.Next())
	}
}
