package transaction

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID
	Description string
	Amount      int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Transactions []Transaction

var NoTransactions = []Transaction{}
var NoTransaction = Transaction{}
