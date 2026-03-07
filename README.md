# go-sqlbuilder

`go-sqlbuilder` is a library that helps you build SQL query strings. It serves to provide a common way to build SQL
query strings regardless of the dialect you're using. It is _NOT_ an ORM.

### Getting Started

`go-sqlbuilder` is easy to use:

```go
import (
	"database/sql"
	"testing"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/cszczepaniak/gotest/assert"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/formatter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/table"
)

db, err := sql.Open(`sqlite3`, `:memory:`)
assert.NoError(t, err)

var myTable = table.Named("MyTable")

b := sqlbuilder.New(formatter.Sqlite{})

// Create a table
_, err = b.CreateTable(myTable).
	Columns(
		column.VarChar("ID", 32).NotNull().PrimaryKey(),
		column.Int("NumberField"),
		column.VarChar("TextField", 255),
	).
	Exec(db)
assert.NoError(t, err)

// Insert some data
_, err = b.InsertInto(myTable).
	Columns("ID", "NumberField", "TextField").
	Values("a", 1, "aa").
	Values("b", 2, "bb").
	Values("c", 3, "cc").
	Exec(db)
assert.NoError(t, err)

// Query your data
row, err := b.SelectFrom(myTable).
	Columns("NumberField", "TextField").
	Where(filter.Equals("NumberField", 3)).
	QueryRow(db) // Or Query
assert.NoError(t, err)

var numField int
var stringField string

err = row.Scan(&numField, &stringField)
assert.NoError(t, err)

assert.Equal(t, numField, 3)
assert.Equal(t, stringField, "cc")

// Update your data
_, err = b.Update(myTable).
	SetFieldTo("NumberField", 123).
	SetFieldTo("TextField", "gotcha").
	Where(filter.Equals("NumberField", 3)).
	Exec(db)
assert.NoError(t, err)

// See the updates
row, err = b.SelectFrom(myTable).
	Columns("NumberField", "TextField").
	Where(filter.Equals("NumberField", 123)).
	QueryRow(db) // Or Query
assert.NoError(t, err)

err = row.Scan(&numField, &stringField)
assert.NoError(t, err)

assert.Equal(t, numField, 123)
assert.Equal(t, stringField, "gotcha")

// Delete your data
res, err := b.DeleteFrom(myTable).
	Where(filter.Greater("NumberField", 10)).
	Exec(db)
assert.NoError(t, err)

n, err := res.RowsAffected()
assert.NoError(t, err)
assert.Equal(t, int(n), 1)
```
