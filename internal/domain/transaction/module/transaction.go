package transaction_module

import (
	common_module "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/module"
	transaction_controller "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/controller"
	transaction_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/repository"
	transaction_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/service"
	"github.com/labstack/echo/v4"
)

type Module struct {
	*common_module.Module
}

func (m *Module) WireController(e *echo.Echo) {
	transactionRepository := transaction_repository.NewPostgresRepository(m.Dependency.Logger, m.Dependency.DatabaseManager)	
	transactionService := transaction_service.New(transactionRepository)
	transactionController := transaction_controller.New(transactionService)
	transactionController.Register(e)
}

func HTTP() common_module.HTTPModule {
	return &Module{
		Module: &common_module.Module{},
	}
}
