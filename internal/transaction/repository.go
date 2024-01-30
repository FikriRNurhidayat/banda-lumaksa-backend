package transaction

import (
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
)

type TransactionRepository repository.Repository[Transaction, TransactionSpecification]
