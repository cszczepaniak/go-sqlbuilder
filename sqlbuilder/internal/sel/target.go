package sel

type Target interface {
	SelectTarget() (string, error)
}

type Table string

func (t Table) SelectTarget() (string, error) {
	return string(t), nil
}
