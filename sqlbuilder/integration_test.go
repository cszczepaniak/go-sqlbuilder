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
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/conflict"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/dialect/mysql"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/dialect/sqlite"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/statement"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func isMySQL(t *testing.T) bool {
	dbChoice := os.Getenv(`TEST_DATABASE`)
	return strings.ToLower(dbChoice) == `mysql`
}

func openSQLiteDatabase(t *testing.T, createTable bool) *sql.DB {
	t.Helper()

	dir, err := os.MkdirTemp(``, ``)
	require.NoError(t, err)

	dataSource := path.Join(dir, `sqlite-database.db`)

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll(dir))
	})

	db, err := sql.Open(`sqlite3`, dataSource)
	require.NoError(t, err)

	if createTable {
		createTestSQLiteTable(t, db)
	}

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

func openMySQLDatabase(t *testing.T, createTable bool) *sql.DB {
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

	if createTable {
		createTestMySQLTable(t, db)
	}

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

func getDatabaseAndBuilder(t *testing.T) (*sql.DB, *sqlbuilder.Builder) {
	if isMySQL(t) {
		t.Log(`--- Using MySQL database for testing ---`)

		db := openMySQLDatabase(t, true)
		b := sqlbuilder.New(mysql.Dialect{})
		return db, b
	}

	t.Log(`--- Using SQLite database for testing ---`)

	db := openSQLiteDatabase(t, true)
	b := sqlbuilder.New(sqlite.Dialect{})
	return db, b
}

func getDatabaseAndBuilderWithoutTable(t *testing.T) (*sql.DB, *sqlbuilder.Builder) {
	if isMySQL(t) {
		t.Log(`--- Using MySQL database for testing ---`)

		db := openMySQLDatabase(t, false)
		b := sqlbuilder.New(mysql.Dialect{})
		return db, b
	}

	t.Log(`--- Using SQLite database for testing ---`)

	db := openSQLiteDatabase(t, false)
	b := sqlbuilder.New(sqlite.Dialect{})
	return db, b
}

func TestMySQLAutoIncrement(t *testing.T) {
	if !isMySQL(t) {
		t.Skip(`test requires MySQL`)
	}

	db := openMySQLDatabase(t, false)
	b := sqlbuilder.New(mysql.Dialect{})

	stmt, err := b.CreateTable(`Test1`).
		Columns(
			column.BigInt(`A`).PrimaryKey().AutoIncrement().Build(),
			column.VarChar(`B`, 20).Build(),
		).
		Build()
	require.NoError(t, err)

	_, err = db.Exec(stmt)
	require.NoError(t, err)

	_, err = b.InsertIntoTable(`Test1`).
		Fields(`B`).
		Values(`AAA`).
		Values(`BBB`).
		Values(`CCC`).
		Exec(db)
	require.NoError(t, err)

	rows, err := b.SelectFromTable(`Test1`).Columns(`A`, `B`).Query(db)
	require.NoError(t, err)

	var (
		aCol int
		bCol string
	)
	assert.True(t, rows.Next())
	require.NoError(t, rows.Scan(&aCol, &bCol))
	assert.Equal(t, 1, aCol)
	assert.Equal(t, `AAA`, bCol)

	assert.True(t, rows.Next())
	require.NoError(t, rows.Scan(&aCol, &bCol))
	assert.Equal(t, 2, aCol)
	assert.Equal(t, `BBB`, bCol)

	assert.True(t, rows.Next())
	require.NoError(t, rows.Scan(&aCol, &bCol))
	assert.Equal(t, 3, aCol)
	assert.Equal(t, `CCC`, bCol)
}

func TestCreateTable(t *testing.T) {
	db, b := getDatabaseAndBuilderWithoutTable(t)
	stmt, err := b.CreateTable(`Test1`).
		Columns(
			column.BigInt(`A`).PrimaryKey().Build(),
			column.BigInt(`B`).Default(123).Build(),
			column.VarChar(`C`, 10).Null().Build(),
		).
		Build()
	require.NoError(t, err)

	_, err = db.Exec(stmt)
	require.NoError(t, err)

	stmt, err = b.CreateTable(`Test1`).
		Columns(
			column.BigInt(`A`).PrimaryKey().Build(),
			column.BigInt(`B`).Default(123).Build(),
			column.VarChar(`C`, 10).Null().Build(),
		).
		Build()
	require.NoError(t, err)

	_, err = db.Exec(stmt)
	// Can't re-create
	require.Error(t, err)

	stmt, err = b.CreateTable(`Test1`).
		IfNotExists().
		Columns(
			column.BigInt(`A`).PrimaryKey().Build(),
			column.BigInt(`B`).Default(123).PrimaryKey().Build(),
			column.VarChar(`C`, 10).Null().Build(),
		).
		Build()
	require.NoError(t, err)

	_, err = db.Exec(stmt)
	// No error with IfNotExists
	require.NoError(t, err)

	_, err = b.InsertIntoTable(`Test1`).
		Fields(`A`, `C`).
		Values(1, `AAA`).
		Values(2, `BBB`).
		Exec(db)
	require.NoError(t, err)

	rows, err := b.SelectFromTable(`Test1`).Columns(`A`, `B`, `C`).Query(db)
	require.NoError(t, err)
	defer rows.Close()

	var (
		aCol int
		bCol int
		cCol string
	)
	assert.True(t, rows.Next())
	require.NoError(t, rows.Scan(&aCol, &bCol, &cCol))
	assert.Equal(t, 1, aCol)
	assert.Equal(t, `AAA`, cCol)

	assert.True(t, rows.Next())
	require.NoError(t, rows.Scan(&aCol, &bCol, &cCol))
	assert.Equal(t, 2, aCol)
	assert.Equal(t, 123, bCol)
	assert.Equal(t, `BBB`, cCol)

	assert.False(t, rows.Next())
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
}

func TestInsertBatches(t *testing.T) {
	db, b := getDatabaseAndBuilder(t)

	execStmts := func(stmts []statement.Statement) {
		for _, stmt := range stmts {
			_, err := db.Exec(stmt.Stmt, stmt.Args...)
			require.NoError(t, err)
		}
	}

	validateTable := func(t *testing.T, exp ...[3]any) {
		rows, err := b.SelectFromTable(`Example`).
			Columns(`ID`, `NumberField`, `TextField`).
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

		assert.Equal(t, i, len(exp), `expected to scan %d rows`, len(exp))

		_, err = b.DeleteFromTable(`Example`).Exec(db)
		require.NoError(t, err)
	}

	stmts, err := b.InsertIntoTable(`Example`).
		Fields(`ID`, `NumberField`, `TextField`).
		Values(`a`, 1, `aa`).
		BuildBatchesOfSize(3)
	require.NoError(t, err)
	assert.Len(t, stmts, 1)

	execStmts(stmts)
	validateTable(t,
		[3]any{`a`, 1, `aa`},
	)

	stmts, err = b.InsertIntoTable(`Example`).
		Fields(`ID`, `NumberField`, `TextField`).
		Values(`a`, 1, `aa`).
		Values(`b`, 2, `bb`).
		Values(`c`, 3, `cc`).
		BuildBatchesOfSize(3)
	require.NoError(t, err)
	assert.Len(t, stmts, 1)

	execStmts(stmts)
	validateTable(t,
		[3]any{`a`, 1, `aa`},
		[3]any{`b`, 2, `bb`},
		[3]any{`c`, 3, `cc`},
	)

	stmts, err = b.InsertIntoTable(`Example`).
		Fields(`ID`, `NumberField`, `TextField`).
		Values(`a`, 1, `aa`).
		Values(`b`, 2, `bb`).
		Values(`c`, 3, `cc`).
		Values(`d`, 4, `dd`).
		BuildBatchesOfSize(3)
	require.NoError(t, err)
	assert.Len(t, stmts, 2)

	execStmts(stmts)
	validateTable(t,
		[3]any{`a`, 1, `aa`},
		[3]any{`b`, 2, `bb`},
		[3]any{`c`, 3, `cc`},
		[3]any{`d`, 4, `dd`},
	)

	stmts, err = b.InsertIntoTable(`Example`).
		Fields(`ID`, `NumberField`, `TextField`).
		Values(`a`, 1, `aa`).
		Values(`b`, 2, `bb`).
		Values(`c`, 3, `cc`).
		Values(`d`, 4, `dd`).
		Values(`e`, 5, `ee`).
		Values(`f`, 6, `ff`).
		Values(`g`, 7, `gg`).
		Values(`h`, 8, `hh`).
		Values(`i`, 9, `ii`).
		BuildBatchesOfSize(3)
	require.NoError(t, err)
	assert.Len(t, stmts, 3)

	execStmts(stmts)
	validateTable(t,
		[3]any{`a`, 1, `aa`},
		[3]any{`b`, 2, `bb`},
		[3]any{`c`, 3, `cc`},
		[3]any{`d`, 4, `dd`},
		[3]any{`e`, 5, `ee`},
		[3]any{`f`, 6, `ff`},
		[3]any{`g`, 7, `gg`},
		[3]any{`h`, 8, `hh`},
		[3]any{`i`, 9, `ii`},
	)
}

func TestConflicts(t *testing.T) {
	db, b := getDatabaseAndBuilder(t)

	validateTable := func(t *testing.T, exp ...[3]any) {
		rows, err := b.SelectFromTable(`Example`).
			Columns(`ID`, `NumberField`, `TextField`).
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

		assert.Equal(t, i, len(exp), `expected to scan %d rows`, len(exp))
	}

	_, err := b.InsertIntoTable(`Example`).
		Fields(`ID`, `NumberField`, `TextField`).
		Values(`a`, 1, `aa`).
		Values(`b`, 2, `bb`).
		Values(`c`, 3, `cc`).
		Values(`d`, 4, `dd`).
		Values(`e`, 5, `ee`).
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
	_, err = b.InsertIntoTable(`Example`).
		Fields(`ID`, `NumberField`, `TextField`).
		Values(`a`, 123, `abc`).
		Values(`f`, 6, `ff`).
		Exec(db)
	require.Error(t, err)

	validateTable(t,
		[3]any{`a`, 1, `aa`},
		[3]any{`b`, 2, `bb`},
		[3]any{`c`, 3, `cc`},
		[3]any{`d`, 4, `dd`},
		[3]any{`e`, 5, `ee`},
	)

	_, err = b.InsertIntoTable(`Example`).
		Fields(`ID`, `NumberField`, `TextField`).
		Values(`a`, 123, `abc`).
		Values(`f`, 6, `ff`).
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

	_, err = b.InsertIntoTable(`Example`).
		Fields(`ID`, `NumberField`, `TextField`).
		Values(`a`, 123, `abc`).
		Values(`f`, 6, `ff`).
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

	_, err = b.InsertIntoTable(`Example`).
		Fields(`ID`, `NumberField`, `TextField`).
		Values(`a`, 1, `def`).
		Values(`f`, 6, `ff`).
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

	res, err := b.InsertIntoTable(`Example`).
		Fields(`ID`, `NumberField`, `TextField`).
		Values(`a`, 1, `aa`).
		Values(`b`, 2, `bb`).
		Values(`c`, 3, `cc`).
		Values(`d`, 4, `dd`).
		Values(`e`, 5, `ee`).
		Exec(db)
	require.NoError(t, err)

	n, err := res.RowsAffected()
	require.NoError(t, err)
	assert.EqualValues(t, 5, n)

	row, err := b.SelectFromTable(`Example`).
		Columns(`NumberField`, `TextField`).
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

	rows, err := b.SelectFromTable(`Example`).
		Columns(`ID`, `NumberField`, `TextField`).
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

	rows, err = b.SelectFromTable(`Example`).
		Columns(`ID`, `NumberField`, `TextField`).
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

	res, err = b.UpdateTable(`Example`).
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

	row, err = b.SelectFromTable(`Example`).
		Columns(`*`).
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

	res, err = b.DeleteFromTable(`Example`).
		Where(filter.Greater(`NumberField`, 3)).
		Exec(tx)
	require.NoError(t, err)

	n, err = res.RowsAffected()
	require.NoError(t, err)
	assert.EqualValues(t, 3, n)

	require.NoError(t, tx.Commit())

	rows, err = b.SelectFromTable(`Example`).
		Columns(`ID`, `NumberField`, `TextField`).
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
