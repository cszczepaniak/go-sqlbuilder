package filter

type Direction int

const (
	Ascending Direction = iota
	Descending
)

type Order struct {
	Column    string
	Direction Direction
}

func OrderDesc(field string) Order {
	return Order{
		Column:    field,
		Direction: Descending,
	}
}

func OrderAsc(field string) Order {
	return Order{
		Column:    field,
		Direction: Ascending,
	}
}
