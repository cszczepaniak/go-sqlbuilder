package ast

type LockKind int

const (
	NoLock LockKind = iota
	SharedLock
	ForUpdateLock
)

type Lock struct {
	Node
	Kind LockKind
}
