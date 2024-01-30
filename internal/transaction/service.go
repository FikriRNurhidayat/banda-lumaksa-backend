package transaction

import (
	"context"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/service"
	"github.com/fikrirnurhidayat/banda-lumaksa/pkg/exists"
	"github.com/google/uuid"
)

type TransactionService interface {
	GetTransaction(ctx context.Context, params *GetTransactionParams) (*GetTransactionResult, error)
	ListTransactions(ctx context.Context, params *ListTransactionsParams) (*ListTransactionsResult, error)
}

type GetTransactionParams struct {
	ID uuid.UUID
}

type GetTransactionResult struct {
	Transaction Transaction
}

type ListTransactionsParams struct {
	DescriptionLike string
	Pagination      service.PaginationParams
}

type ListTransactionsResult struct {
	Pagination   service.PaginationResult
	Transactions []Transaction
}

type TransactionServiceImpl struct {
	transactionRepository TransactionRepository
}

// GetTranscation implements TransactionService.
func (s *TransactionServiceImpl) GetTransaction(ctx context.Context, params *GetTransactionParams) (*GetTransactionResult, error) {
	transaction, err := s.transactionRepository.Get(ctx, WithID(params.ID))
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	if transaction == NoTransaction {
		return nil, ErrTransactionNotFound
	}

	return &GetTransactionResult{
		Transaction: transaction,
	}, nil
}

func (s *TransactionServiceImpl) ListTransactions(ctx context.Context, params *ListTransactionsParams) (*ListTransactionsResult, error) {
	filters := []TransactionSpecification{}

	if exists.String(params.DescriptionLike) {
		filters = append(filters, DescriptionLike(params.DescriptionLike))
	}

	params.Pagination = params.Pagination.Normalize()

	transactions, err := s.transactionRepository.List(ctx, repository.ListArgs[TransactionSpecification]{
		Filters: filters,
		Limit:   params.Pagination.Limit(),
		Offset:  params.Pagination.Offset(),
	})
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	size, err := s.transactionRepository.Size(ctx, filters...)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return &ListTransactionsResult{
		Pagination:   service.NewPaginationResult(params.Pagination, size),
		Transactions: transactions,
	}, nil
}

func NewTransactionService(transactionRepository TransactionRepository) TransactionService {
	return &TransactionServiceImpl{
		transactionRepository: transactionRepository,
	}
}
