# go-sqlbuilder

`go-sqlbuilder` is a library that helps you build SQL query strings. It serves to provide a common way to build SQL
query strings regardless of the dialect you're using. It is _NOT_ an ORM.

### Getting Started

`go-sqlbuilder` is easy to use:

```go
import (
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/dialect/sqlite"
	"github.com/cszczepaniak/go-sqlbuilder/sqlbuilder/filter"
)

var db *sql.DB

b := sqlbuilder.New(sqlite.Dialect{})

// Insert some data
_, err := b.Insert("MyTable").
    Fields("ID", "NumberField", "TextField").
    WithRecord("a", 1, "aa").
    WithRecord("b", 2, "bb").
    WithRecord("c", 3, "cc").
    Exec(db)

// Query your data
row, err := b.Select("MyTable").
	Fields("NumberField", "TextField").
	Where(filter.Equals("NumberField", 3)).
	QueryRow(db) // Or Query

// Update your data
_, err = b.Update("MyTable").
	SetFieldTo("NumberField", 123).
	SetFieldTo("TextField", "gotcha").
	WhereAll(
		filter.NotEquals("TextField", "bb"),
		filter.LessOrEqual("NumberField", 2),
	).
	Exec(db)

// Delete your data
_, err = b.Delete("MyTable").
	Where(filter.Greater("NumberField", 3)).
	Exec(db)

```
