package sqlbuilder

type Dialect interface {
	selectDialect
	deleteDialect
	updateDialect
	insertDialect
}

type QueryBuilder struct {
	d        Dialect
	database string
}

func New(d Dialect) *QueryBuilder {
	return &QueryBuilder{
		d: d,
	}
}

func (b *QueryBuilder) SetDatabase(db string) *QueryBuilder {
	b.database = db
	return b
}

func (b *QueryBuilder) qualifiedTableName(table string) string {
	if b.database != `` {
		return b.database + `.` + table
	}
	return table
}

func (b *QueryBuilder) Select(table string) *SelectBuilder {
	return newSelectBuilder(b.d, b.qualifiedTableName(table))
}

func (b *QueryBuilder) Delete(table string) *DeleteBuilder {
	return newDeleteBuilder(b.d, b.qualifiedTableName(table))
}

func (b *QueryBuilder) Update(table string) *UpdateBuilder {
	return newUpdateBuilder(b.d, b.qualifiedTableName(table))
}

func (b *QueryBuilder) Insert(table string) *InsertBuilder {
	return newInsertBuilder(b.d, b.qualifiedTableName(table))
}

type TableQueryBuilder struct {
	b     *QueryBuilder
	table string
}

func (b *QueryBuilder) ForTable(table string) *TableQueryBuilder {
	return &TableQueryBuilder{
		b:     b,
		table: table,
	}
}

func (b *TableQueryBuilder) Select() *SelectBuilder {
	return newSelectBuilder(b.b.d, b.table)
}

func (b *TableQueryBuilder) Delete() *DeleteBuilder {
	return newDeleteBuilder(b.b.d, b.table)
}

func (b *TableQueryBuilder) Update() *UpdateBuilder {
	return newUpdateBuilder(b.b.d, b.table)
}

func (b *TableQueryBuilder) Insert() *InsertBuilder {
	return newInsertBuilder(b.b.d, b.table)
}
