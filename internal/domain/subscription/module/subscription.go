package subscription_module

import (
	common_module "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/module"
	subscription_command "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/command"
	subscription_controller "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/controller"
	subscription_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/repository"
	subscription_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/service"
	transaction_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/repository"
	"github.com/labstack/echo/v4"
)

type Module struct {
	*common_module.Module
}

func (m *Module) WireController(e *echo.Echo) {
	subscriptionRepository := subscription_repository.NewPostgresRepository(m.Dependency.Logger, m.Dependency.DatabaseManager)
	transactionRepository := transaction_repository.NewPostgresRepository(m.Dependency.Logger, m.Dependency.DatabaseManager)
	subscriptionService := subscription_service.New(m.Dependency.Logger, subscriptionRepository, transactionRepository, m.Dependency.TransactionManager)
	subscriptionController := subscription_controller.New(m.Dependency.Logger, subscriptionService)
	subscriptionController.Register(e)
}

func (m *Module) WireCommand() subscription_command.SubscriptionCommand {
	subscriptionRepository := subscription_repository.NewPostgresRepository(m.Dependency.Logger, m.Dependency.DatabaseManager)
	transactionRepository := transaction_repository.NewPostgresRepository(m.Dependency.Logger, m.Dependency.DatabaseManager)
	subscriptionService := subscription_service.New(m.Dependency.Logger, subscriptionRepository, transactionRepository, m.Dependency.TransactionManager)
	subscriptionCommand := subscription_command.New(subscriptionService)

	return subscriptionCommand
}

func HTTP() common_module.HTTPModule {
	return &Module{
		Module: &common_module.Module{},
	}
}
