package ast

func Visit(n Node, fn func(n Node) bool) {
	n.AcceptVisitor(fn)
}

func GetArgs(n Node) []any {
	var args []any

	Visit(n, func(n Node) bool {
		ph, ok := n.(*PlaceholderLiteral)
		if !ok {
			// Keep traversing the tree
			return true
		}
		args = append(args, ph.For)
		return false
	})

	return args
}
