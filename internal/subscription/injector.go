package subscription

import (
	"database/sql"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/transaction"
)

func InitializeHandlers(db *sql.DB) SubscrpitionController {
	databaseManager := manager.NewDatabaseManager(db)
	transactionManager := manager.NewTransactionManager(db)
	subscriptionRepository := NewPostgresSubscriptionRepository(databaseManager)
	transactionRepository := transaction.NewPostgresTransactionRepository(databaseManager)
	subscriptionService := NewSubscriptionService(subscriptionRepository, transactionRepository, transactionManager)
	subscriptionController := NewSubscriptionController(subscriptionService)

	return subscriptionController
}

func InitializeCommands(db *sql.DB) SubscriptionCommand {
	databaseManager := manager.NewDatabaseManager(db)
	transactionManager := manager.NewTransactionManager(db)
	subscriptionRepository := NewPostgresSubscriptionRepository(databaseManager)
	transactionRepository := transaction.NewPostgresTransactionRepository(databaseManager)
	subscriptionService := NewSubscriptionService(subscriptionRepository, transactionRepository, transactionManager)
	subscriptionCommand := NewSubscriptionCommand(subscriptionService)

	return subscriptionCommand
}
