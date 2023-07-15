package sqlbuilder

type TableQueryBuilder struct {
	table string
}

func NewQueryBuilder(table string) *TableQueryBuilder {
	return &TableQueryBuilder{
		table: table,
	}
}

func (tb *TableQueryBuilder) Select(d selectDialect) *SelectBuilder {
	return newSelectBuilder(d, tb.table)
}

func (tb *TableQueryBuilder) Delete(d deleteDialect) *DeleteBuilder {
	return newDeleteBuilder(d, tb.table)
}

func (tb *TableQueryBuilder) Update(d updateDialect) *UpdateBuilder {
	return newUpdateBuilder(d, tb.table)
}

func (tb *TableQueryBuilder) Insert(d insertDialect) *InsertBuilder {
	return newInsertBuilder(d, tb.table)
}
