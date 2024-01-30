package transaction

import (
	"database/sql"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"
)

func InitializeHandlers(db *sql.DB) TransactionController {
	databaseManager := manager.NewDatabaseManager(db)
	transactionRepository := NewPostgresTransactionRepository(databaseManager)
	transactionService := NewTransactionService(transactionRepository)
	return NewTransactionController(transactionService)
}
