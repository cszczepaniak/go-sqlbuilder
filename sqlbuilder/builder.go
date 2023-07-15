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
