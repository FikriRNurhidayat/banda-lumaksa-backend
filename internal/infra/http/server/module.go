package http_server

import (
	subscription_controller "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/controller"
	subscription_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/repository"
	subscription_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/service"
	transaction_controller "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/controller"
	transaction_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/repository"
	transaction_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/service"
)

type Dependency struct {
	TransactionRepository  transaction_repository.TransactionRepository
	TransactionService     transaction_service.TransactionService
	TransactionController  transaction_controller.TransactionController
	SubscriptionRepository subscription_repository.SubscriptionRepository
	SubscriptionService    subscription_service.SubscriptionService
	SubscriptionController subscription_controller.SubscriptionController
}

func (s *Server) Bootstrap() (err error) {
	s.Dependency = &Dependency{}

	s.Dependency.SubscriptionRepository, err = subscription_repository.NewPostgresRepository(s.RootDependency.Logger, s.RootDependency.DatabaseManager)
	if err != nil {
		return err
	}

	s.Dependency.TransactionRepository, err = transaction_repository.NewPostgresRepository(s.RootDependency.Logger, s.RootDependency.DatabaseManager)
	if err != nil {
		return err
	}

	s.Dependency.TransactionService = transaction_service.New(s.Dependency.TransactionRepository)
	s.Dependency.SubscriptionService = subscription_service.New(s.RootDependency.Logger, s.Dependency.SubscriptionRepository, s.Dependency.TransactionRepository, s.RootDependency.TransactionManager)

	s.Dependency.SubscriptionController = subscription_controller.New(s.Logger, s.Dependency.SubscriptionService)
	s.Dependency.TransactionController = transaction_controller.New(s.Dependency.TransactionService)

	s.Dependency.SubscriptionController.Register(s.Echo)
	s.Dependency.TransactionController.Register(s.Echo)

	return nil
}
