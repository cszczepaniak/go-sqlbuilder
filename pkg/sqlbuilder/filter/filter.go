package filter

type Filter interface {
	Args() []any
}

type AllFilter struct {
	Filters []Filter
}

func All(fs ...Filter) AllFilter {
	return AllFilter{
		Filters: fs,
	}
}

func (f AllFilter) Args() []any {
	var args []any
	for _, ff := range f.Filters {
		args = append(args, ff.Args()...)
	}
	return args
}

type AnyFilter struct {
	Filters []Filter
}

func Any(fs ...Filter) AnyFilter {
	return AnyFilter{
		Filters: fs,
	}
}

func (f AnyFilter) Args() []any {
	var args []any
	for _, ff := range f.Filters {
		args = append(args, ff.Args()...)
	}
	return args
}

type EqualsFilter struct {
	Column string
	Value  any
}

func Equals(column string, val any) EqualsFilter {
	return EqualsFilter{
		Column: column,
		Value:  val,
	}
}

func (f EqualsFilter) Args() []any {
	return []any{f.Value}
}

type NotEqualsFilter struct {
	Column string
	Value  any
}

func NotEquals(column string, val any) NotEqualsFilter {
	return NotEqualsFilter{
		Column: column,
		Value:  val,
	}
}

func (f NotEqualsFilter) Args() []any {
	return []any{f.Value}
}

type GreaterFilter struct {
	Column string
	Value  any
}

func Greater(column string, val any) GreaterFilter {
	return GreaterFilter{
		Column: column,
		Value:  val,
	}
}

func (f GreaterFilter) Args() []any {
	return []any{f.Value}
}

type GreaterOrEqualFilter struct {
	Column string
	Value  any
}

func GreaterOrEqual(column string, val any) GreaterOrEqualFilter {
	return GreaterOrEqualFilter{
		Column: column,
		Value:  val,
	}
}

func (f GreaterOrEqualFilter) Args() []any {
	return []any{f.Value}
}

type LessFilter struct {
	Column string
	Value  any
}

func Less(column string, val any) LessFilter {
	return LessFilter{
		Column: column,
		Value:  val,
	}
}

func (f LessFilter) Args() []any {
	return []any{f.Value}
}

type LessOrEqualFilter struct {
	Column string
	Value  any
}

func LessOrEqual(column string, val any) LessOrEqualFilter {
	return LessOrEqualFilter{
		Column: column,
		Value:  val,
	}
}

func (f LessOrEqualFilter) Args() []any {
	return []any{f.Value}
}

type InFilter struct {
	Column string
	Values []any
}

func In(column string, vals ...any) InFilter {
	return InFilter{
		Column: column,
		Values: vals,
	}
}

func (f InFilter) Args() []any {
	return f.Values
}
