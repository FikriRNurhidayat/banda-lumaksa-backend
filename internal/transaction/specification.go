package transaction

import (
	"strings"

	"github.com/google/uuid"
)

type TransactionSpecification interface {
	Call(transaction Transaction) bool
}

type DescriptionLikeSpecification struct {
	Like string
}

func (spec DescriptionLikeSpecification) Call(transaction Transaction) bool {
	return strings.Contains(strings.ToLower(transaction.Description), strings.ToLower(spec.Like))
}

func DescriptionLike(like string) TransactionSpecification {
	return DescriptionLikeSpecification{
		Like: like,
	}
}

type WithIDSpecification struct {
	ID uuid.UUID
}

func (spec WithIDSpecification) Call(transaction Transaction) bool {
	return spec.ID == transaction.ID
}

func WithID(id uuid.UUID) TransactionSpecification {
	return WithIDSpecification{
		ID: id,
	}
}
