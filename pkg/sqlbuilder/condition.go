package sqlbuilder

import "github.com/cszczepaniak/go-sqlbuilder/pkg/sqlbuilder/filter"

type conditioner interface {
	Condition(f filter.Filter) (string, error)
}

func getCondition(c conditioner, f filter.Filter) (string, []any, error) {
	if f != nil {
		cond, err := c.Condition(f)
		if err != nil {
			return ``, nil, err
		}
		return cond, f.Args(), nil
	}
	return ``, nil, nil
}
