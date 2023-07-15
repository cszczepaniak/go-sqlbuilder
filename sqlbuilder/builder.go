package sqlbuilder

type QueryBuilder struct {
	table string
}

func New(table string) *QueryBuilder {
	return &QueryBuilder{
		table: table,
	}
}

func (tb *QueryBuilder) Select(d selectDialect) *SelectBuilder {
	return newSelectBuilder(d, tb.table)
}

func (tb *QueryBuilder) Delete(d deleteDialect) *DeleteBuilder {
	return newDeleteBuilder(d, tb.table)
}

func (tb *QueryBuilder) Update(d updateDialect) *UpdateBuilder {
	return newUpdateBuilder(d, tb.table)
}

func (tb *QueryBuilder) Insert(d insertDialect) *InsertBuilder {
	return newInsertBuilder(d, tb.table)
}
