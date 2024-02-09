package transaction_repository

import (
	common_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"

	transaction_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/entity"
	transaction_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/specification"
)

type TransactionRepository common_repository.Repository[transaction_entity.Transaction, transaction_specification.TransactionSpecification]
