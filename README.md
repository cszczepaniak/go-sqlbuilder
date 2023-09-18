# go-sqlbuilder

`go-sqlbuilder` is a library that helps you build SQL query strings. It serves to provide a common way to build SQL
query strings regardless of the dialect you're using. It is _NOT_ an ORM.

### Getting Started

`go-sqlbuilder` is easy to use:

```go
import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/column"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/dialect/sqlite"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/formatter"
)

db, err := sql.Open(`sqlite3`, `:memory:`)
require.NoError(t, err)

b := sqlbuilder.New(sqlite.Dialect{}, formatter.Sqlite{})

// Create a table
stmt, err := b.CreateTable("MyTable").
	Columns(
		column.VarChar("ID", 32).NotNull().PrimaryKey().Build(),
		column.Int("NumberField").Build(),
		column.VarChar("TextField", 255).Build(),
	).
	Build()
require.NoError(t, err)

_, err = db.Exec(stmt)
require.NoError(t, err)

// Insert some data
_, err = b.InsertIntoTable("MyTable").
	Fields("ID", "NumberField", "TextField").
	Values("a", 1, "aa").
	Values("b", 2, "bb").
	Values("c", 3, "cc").
	Exec(db)
require.NoError(t, err)

// Query your data
row, err := b.SelectFromTable("MyTable").
	Columns("NumberField", "TextField").
	Where(filter.Equals("NumberField", 3)).
	QueryRow(db) // Or Query
require.NoError(t, err)

var numField int
var stringField string

err = row.Scan(&numField, &stringField)
require.NoError(t, err)

assert.Equal(t, 3, numField)
assert.Equal(t, "cc", stringField)

// Update your data
_, err = b.UpdateTable("MyTable").
	SetFieldTo("NumberField", 123).
	SetFieldTo("TextField", "gotcha").
	Where(filter.Equals("NumberField", 3)).
	Exec(db)
require.NoError(t, err)

// See the updates
row, err = b.SelectFromTable("MyTable").
	Columns("NumberField", "TextField").
	Where(filter.Equals("NumberField", 123)).
	QueryRow(db) // Or Query
require.NoError(t, err)

err = row.Scan(&numField, &stringField)
require.NoError(t, err)

assert.Equal(t, 123, numField)
assert.Equal(t, "gotcha", stringField)

// Delete your data
res, err := b.DeleteFromTable("MyTable").
	Where(filter.Greater("NumberField", 10)).
	Exec(db)
require.NoError(t, err)

n, err := res.RowsAffected()
require.NoError(t, err)
assert.EqualValues(t, 1, n)
```
