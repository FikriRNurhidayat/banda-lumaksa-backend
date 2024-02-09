package subscription_initializer

import (
	"database/sql"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/logger"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/transaction"

	subscription_command "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/command"
	subscription_controller "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/controller"
	subscription_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/repository"
	subscription_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/service"
)

func InitializeController(db *sql.DB, logger logger.Logger) subscription_controller.SubscriptionController {
	databaseManager := manager.NewDatabaseManager(db)
	transactionManager := manager.NewTransactionManager(db)
	subscriptionRepository := subscription_repository.NewPostgresRepository(databaseManager)
	transactionRepository := transaction.NewPostgresTransactionRepository(databaseManager)
	subscriptionService := subscription_service.New(logger, subscriptionRepository, transactionRepository, transactionManager)
	subscriptionController := subscription_controller.New(logger, subscriptionService)

	return subscriptionController
}

func InitializeCommand(db *sql.DB, logger logger.Logger) subscription_command.SubscriptionCommand {
	databaseManager := manager.NewDatabaseManager(db)
	transactionManager := manager.NewTransactionManager(db)
	subscriptionRepository := subscription_repository.NewPostgresRepository(databaseManager)
	transactionRepository := transaction.NewPostgresTransactionRepository(databaseManager)
	subscriptionService := subscription_service.New(logger, subscriptionRepository, transactionRepository, transactionManager)
	subscriptionCommand := subscription_command.New(subscriptionService)

	return subscriptionCommand
}
