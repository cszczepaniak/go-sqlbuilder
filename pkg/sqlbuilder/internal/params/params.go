package params

import "strings"

func Group(n int) string {
	params := strings.Repeat(`?,`, n)
	params = strings.TrimRight(params, `,`)
	return `(` + params + `)`
}

func Groups(nGroups, nPerGroup int) string {
	group := Group(nPerGroup)
	groups := strings.Repeat(group+`,`, nGroups)
	return strings.TrimRight(groups, `,`)
}
