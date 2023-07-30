package functions

type Count struct {
	Field    string
	Distinct bool
}

func CountAll() Count {
	return Count{}
}

func CountField(f string) Count {
	return Count{
		Field:    f,
		Distinct: false,
	}
}

func CountDistinct(f string) Count {
	return Count{
		Field:    f,
		Distinct: true,
	}
}

func (c Count) Args() []any {
	return nil
}

func (c Count) All() bool {
	return c.Field == ``
}
