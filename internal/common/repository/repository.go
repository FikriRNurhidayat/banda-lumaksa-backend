package repository

import (
	"context"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/specification"
)

type ListArgs[T any] struct {
	Filters []T
	Sort    specification.Specification
	Limit   specification.Specification
	Offset  specification.Specification
}

type Repository[Entity any, Specification any] interface {
	Save(context.Context, Entity) error
	Get(context.Context, ...Specification) (Entity, error)
	Delete(context.Context, ...Specification) error
	List(context.Context, ListArgs[Specification]) ([]Entity, error)
	Each(context.Context, ListArgs[Specification]) (Iterator[Entity], error)
	Size(context.Context, ...Specification) (uint32, error)
}

type Iterator[Entity any] interface {
	Next() bool
	Entry() (Entity, error)
}
