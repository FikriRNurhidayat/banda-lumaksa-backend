package common_module

import (
	"database/sql"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/logger"
	database_manager "github.com/fikrirnurhidayat/banda-lumaksa/internal/manager/database"
	transaction_manager "github.com/fikrirnurhidayat/banda-lumaksa/internal/manager/transaction"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Dependency *RootDependency
}

func (m *Module) Wire(dependency *RootDependency) {
	m.Dependency = dependency
}

type RootDependency struct {
	DatabaseManager    database_manager.DatabaseManager
	TransactionManager transaction_manager.TransactionManager
	Logger             logger.Logger
}

type HTTPModule interface {
	Wire(*RootDependency)
	WireController(e *echo.Echo)
}

type HTTPModules []HTTPModule

type CLIModule interface {
	Wire(*RootDependency)
	WireCommand()
}

type CLIModules []CLIModule

func New(db *sql.DB, logger logger.Logger) *RootDependency {
	databaseManager := database_manager.New(db)
	transactionManager := transaction_manager.New(db)

	return &RootDependency{
		DatabaseManager:    databaseManager,
		TransactionManager: transactionManager,
		Logger:             logger,
	}
}
