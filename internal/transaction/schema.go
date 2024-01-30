package transaction

import (
	"time"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/schema"
	"github.com/google/uuid"
)

type TransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Amount      int32     `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TransactionsResponse []TransactionResponse

type ListTransactionsResponse struct {
	schema.PaginationResponse
	Transactions TransactionsResponse `json:"transactions"`
}

type GetTransactionResponse struct {
	Transaction TransactionResponse `json:"transaction"`
}

func NewTransactionResponse(transaction Transaction) TransactionResponse {
	return TransactionResponse{
		ID:          transaction.ID,
		Description: transaction.Description,
		Amount:      transaction.Amount,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}
}

func NewTransactionsResponse(transactions Transactions) TransactionsResponse {
	transactionsResponse := TransactionsResponse{}

	for _, s := range transactions {
		transactionsResponse = append(transactionsResponse, NewTransactionResponse(s))
	}

	return transactionsResponse
}
