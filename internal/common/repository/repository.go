package common_repository

import (
	"context"

	common_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/specification"
)

type ListArgs[T any] struct {
	Filters []T
	Sort    common_specification.Specification
	Limit   common_specification.Specification
	Offset  common_specification.Specification
}

type Repository[Entity any, Specification any] interface {
	Save(context.Context, Entity) error
	Get(context.Context, ...Specification) (Entity, error)
	Exist(context.Context, ...Specification) (bool, error)
	Delete(context.Context, ...Specification) error
	List(context.Context, ListArgs[Specification]) ([]Entity, error)
	Each(context.Context, ListArgs[Specification]) (Iterator[Entity], error)
	Size(context.Context, ...Specification) (uint32, error)
}

type Iterator[Entity any] interface {
	Next() bool
	Current() (Entity, error)
}
