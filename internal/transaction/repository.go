package transaction

import (
	"context"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	Save(context.Context, Transaction) error
	Get(context.Context, uuid.UUID) (Transaction, error)
	Delete(context.Context, uuid.UUID) error
	List(context.Context, ...TransactionSpecification) (Transactions, error)
	Size(context.Context, ...TransactionSpecification) (uint32, error)
}
