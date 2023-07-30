package limit

type Limiter interface {
	Limit() (string, error)
}

// LimitBuilder is used to construct SQL limits. This builder is meant to be embedded in any other builder which can
// have a limit (e.g. Select)
//
// Note: normally, this name shouldn't stutter with the package name, but we may want to embed several things called
// "Builder" in other builders, so we have to disambiguate.
type LimitBuilder[T any] struct {
	parent T

	limit *int
}

func NewBuilder[T any](parent T) *LimitBuilder[T] {
	return &LimitBuilder[T]{
		parent: parent,
	}
}

func (b *LimitBuilder[T]) Limit(limit int) T {
	b.limit = &limit
	return b.parent
}

func (b *LimitBuilder[T]) SQLAndArgs(l Limiter) (string, []any, error) {
	if b.limit == nil {
		return ``, nil, nil
	}

	lim, err := l.Limit()
	if err != nil {
		return ``, nil, err
	}
	return lim, []any{*b.limit}, nil
}
