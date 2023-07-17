package sqlbuilder_test

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/conflict"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/dialect/mysql"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/dialect/sqlite"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	_ "github.com/go-sql-driver/mysql"
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

	createTestSQLiteTable(t, db)

	t.Cleanup(func() {
		assert.NoError(t, db.Close())
	})

	return db
}

func createTestSQLiteTable(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`CREATE TABLE Example (
		ID TEXT NOT NULL PRIMARY KEY,
		NumberField INT,
		TextField TEXT
	)`)
	require.NoError(t, err)
}

func openMySQLDatabase(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/")
	require.NoError(t, err)

	buff := make([]byte, 0, 16)
	buff = binary.LittleEndian.AppendUint64(buff, uint64(rand.Int63()))
	buff = binary.LittleEndian.AppendUint64(buff, uint64(rand.Int63()))

	dbName := fmt.Sprintf(`test_%x`, buff)

	_, err = db.Exec(`CREATE DATABASE ` + dbName)
	require.NoError(t, err)
	_, err = db.Exec(`USE ` + dbName)
	require.NoError(t, err)

	createTestMySQLTable(t, db)

	t.Cleanup(func() {
		_, err = db.Exec(`DROP DATABASE ` + dbName)
		assert.NoError(t, err)
	})

	return db
}

func createTestMySQLTable(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`CREATE TABLE Example (
		ID VARCHAR(255) NOT NULL PRIMARY KEY,
		NumberField INT,
		TextField TEXT
	)`)
	require.NoError(t, err)
}

func getDatabaseAndBuilder(t *testing.T) (*sql.DB, *sqlbuilder.TableQueryBuilder) {
	dbChoice := os.Getenv(`TEST_DATABASE`)
	if strings.ToLower(dbChoice) == `mysql` {
		t.Log(`--- Using MySQL database for testing ---`)

		db := openMySQLDatabase(t)
		b := sqlbuilder.New(mysql.Dialect{}).ForTable(`Example`)
		return db, b
	}

	t.Log(`--- Using SQLite database for testing ---`)

	db := openSQLiteDatabase(t)
	b := sqlbuilder.New(sqlite.Dialect{}).ForTable(`Example`)
	return db, b
}

func TestConflicts(t *testing.T) {
	db, b := getDatabaseAndBuilder(t)

	validateTable := func(t *testing.T, exp ...[3]any) {
		rows, err := b.Select().
			Fields(`ID`, `NumberField`, `TextField`).
			OrderBy(filter.OrderAsc(`ID`)).
			Query(db)
		require.NoError(t, err)
		defer rows.Close()

		var (
			id     string
			number int
			text   string
		)

		i := 0
		for rows.Next() {
			require.NoError(t, rows.Scan(&id, &number, &text))
			require.LessOrEqual(t, i, len(exp), `expected had additional elements`)

			assert.Equal(t, exp[i][0], id)
			assert.Equal(t, exp[i][1], number)
			assert.Equal(t, exp[i][2], text)
			i++
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}

	_, err := b.Insert().
		Fields(`ID`, `NumberField`, `TextField`).
		WithRecord(`a`, 1, `aa`).
		WithRecord(`b`, 2, `bb`).
		WithRecord(`c`, 3, `cc`).
		WithRecord(`d`, 4, `dd`).
		WithRecord(`e`, 5, `ee`).
		Exec(db)
	require.NoError(t, err)

	validateTable(t,
		[3]any{`a`, 1, `aa`},
		[3]any{`b`, 2, `bb`},
		[3]any{`c`, 3, `cc`},
		[3]any{`d`, 4, `dd`},
		[3]any{`e`, 5, `ee`},
	)

	// Error because of key conflict
	_, err = b.Insert().
		Fields(`ID`, `NumberField`, `TextField`).
		WithRecord(`a`, 123, `abc`).
		WithRecord(`f`, 6, `ff`).
		Exec(db)
	require.Error(t, err)

	validateTable(t,
		[3]any{`a`, 1, `aa`},
		[3]any{`b`, 2, `bb`},
		[3]any{`c`, 3, `cc`},
		[3]any{`d`, 4, `dd`},
		[3]any{`e`, 5, `ee`},
	)

	_, err = b.Insert().
		Fields(`ID`, `NumberField`, `TextField`).
		WithRecord(`a`, 123, `abc`).
		WithRecord(`f`, 6, `ff`).
		IgnoreConflicts(conflict.NewKey(`ID`)).
		Exec(db)
	require.NoError(t, err)

	validateTable(t,
		[3]any{`a`, 1, `aa`},
		[3]any{`b`, 2, `bb`},
		[3]any{`c`, 3, `cc`},
		[3]any{`d`, 4, `dd`},
		[3]any{`e`, 5, `ee`},
		[3]any{`f`, 6, `ff`},
	)

	_, err = b.Insert().
		Fields(`ID`, `NumberField`, `TextField`).
		WithRecord(`a`, 123, `abc`).
		WithRecord(`f`, 6, `ff`).
		OverwriteConflicts(conflict.NewKey(`ID`)).
		Exec(db)
	require.NoError(t, err)

	validateTable(t,
		[3]any{`a`, 123, `abc`},
		[3]any{`b`, 2, `bb`},
		[3]any{`c`, 3, `cc`},
		[3]any{`d`, 4, `dd`},
		[3]any{`e`, 5, `ee`},
		[3]any{`f`, 6, `ff`},
	)

	_, err = b.Insert().
		Fields(`ID`, `NumberField`, `TextField`).
		WithRecord(`a`, 1, `def`).
		WithRecord(`f`, 6, `ff`).
		OnConflict(
			conflict.NewKey(`ID`),
			conflict.Ignore(`NumberField`),
			conflict.Overwrite(`TextField`),
		).
		Exec(db)
	require.NoError(t, err)

	validateTable(t,
		[3]any{`a`, 123, `def`},
		[3]any{`b`, 2, `bb`},
		[3]any{`c`, 3, `cc`},
		[3]any{`d`, 4, `dd`},
		[3]any{`e`, 5, `ee`},
		[3]any{`f`, 6, `ff`},
	)
}

func TestBasicFunction(t *testing.T) {
	db, b := getDatabaseAndBuilder(t)

	res, err := b.Insert().
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

	row, err := b.Select().
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

	rows, err := b.Select().
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

	rows, err = b.Select().
		Fields(`ID`, `NumberField`, `TextField`).
		Where(filter.In(`TextField`, `bb`, `dd`)).
		OrderBy(filter.OrderDesc(`TextField`)).
		Limit(1).
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
		assert.Equal(t, `d`, id)
		assert.Equal(t, 4, numberField)
		assert.Equal(t, `dd`, textField)

		assert.False(t, rows.Next())
	}

	res, err = b.Update().
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

	row, err = b.Select().
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

	res, err = b.Delete().
		Where(filter.Greater(`NumberField`, 3)).
		Exec(tx)
	require.NoError(t, err)

	n, err = res.RowsAffected()
	require.NoError(t, err)
	assert.EqualValues(t, 3, n)

	require.NoError(t, tx.Commit())

	rows, err = b.Select().
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
