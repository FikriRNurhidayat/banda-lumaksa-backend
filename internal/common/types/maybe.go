package types

type Maybe[T any] struct {
	Present bool
	Value   T
}