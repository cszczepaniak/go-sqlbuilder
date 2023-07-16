package conflict

type Behavior interface {
	Field() string
}

type IgnoreBehavior struct {
	field string
}

func Ignore(field string) IgnoreBehavior {
	return IgnoreBehavior{
		field: field,
	}
}

func (b IgnoreBehavior) Field() string { return b.field }

type OverwriteBehavior struct {
	field string
}

func Overwrite(field string) OverwriteBehavior {
	return OverwriteBehavior{
		field: field,
	}
}

func (b OverwriteBehavior) Field() string { return b.field }
