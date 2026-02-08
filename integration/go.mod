module github.com/cszczepaniak/go-sqlbuilder/integration

go 1.25.7

require (
	github.com/cszczepaniak/go-sqlbuilder v0.0.8
	github.com/go-sql-driver/mysql v1.9.3
	github.com/ncruces/go-sqlite3 v0.30.5
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/ncruces/julianday v1.0.0 // indirect
	github.com/tetratelabs/wazero v1.11.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
)

replace github.com/cszczepaniak/go-sqlbuilder => ../
