package sqlbuilder

type limiter interface {
	Limit() (string, error)
}

func getLimit(c limiter, l *int) (string, []any, error) {
	if l != nil {
		lim, err := c.Limit()
		if err != nil {
			return ``, nil, err
		}
		return lim, []any{*l}, nil
	}
	return ``, nil, nil
}
