package transaction_service

import (
	"context"

	common_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
	common_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/service"

	transaction_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/entity"
	transaction_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/errors"
	transaction_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/repository"
	transaction_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/specification"

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
	Transaction transaction_entity.Transaction
}

type ListTransactionsParams struct {
	DescriptionLike string
	Pagination      common_service.PaginationParams
}

type ListTransactionsResult struct {
	Pagination   common_service.PaginationResult
	Transactions []transaction_entity.Transaction
}

type TransactionServiceImpl struct {
	transactionRepository transaction_repository.TransactionRepository
}

// GetTranscation implements TransactionService.
func (s *TransactionServiceImpl) GetTransaction(ctx context.Context, params *GetTransactionParams) (*GetTransactionResult, error) {
	transaction, err := s.transactionRepository.Get(ctx, transaction_specification.WithID(params.ID))
	if err != nil {
		return nil, err
	}

	if transaction == transaction_entity.NoTransaction {
		return nil, transaction_errors.ErrTransactionNotFound
	}

	return &GetTransactionResult{
		Transaction: transaction,
	}, nil
}

func (s *TransactionServiceImpl) ListTransactions(ctx context.Context, params *ListTransactionsParams) (*ListTransactionsResult, error) {
	filters := []transaction_specification.TransactionSpecification{}

	if exists.String(params.DescriptionLike) {
		filters = append(filters, transaction_specification.DescriptionLike(params.DescriptionLike))
	}

	params.Pagination = params.Pagination.Normalize()

	transactions, err := s.transactionRepository.List(ctx, common_repository.ListArgs[transaction_specification.TransactionSpecification]{
		Filters: filters,
		Limit:   params.Pagination.Limit(),
		Offset:  params.Pagination.Offset(),
	})
	if err != nil {
		return nil, err
	}

	size, err := s.transactionRepository.Size(ctx, filters...)
	if err != nil {
		return nil, err
	}

	return &ListTransactionsResult{
		Pagination:   common_service.NewPaginationResult(params.Pagination, size),
		Transactions: transactions,
	}, nil
}

func New(transactionRepository transaction_repository.TransactionRepository) TransactionService {
	return &TransactionServiceImpl{
		transactionRepository: transactionRepository,
	}
}
