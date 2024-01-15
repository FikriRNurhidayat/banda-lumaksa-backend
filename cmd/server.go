package main

import (
	"log"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/server"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription"
)

var (
	subscriptionRepository    subscription.Repository
	createSubscriptionUseCase subscription.CreateSubscriptionUseCase
	cancelSubscriptionUseCase subscription.CancelSubscriptionUseCase
	getSubscriptionUseCase    subscription.GetSubscriptionUseCase
	listSubscriptionsUseCase  subscription.ListSubscriptionsUseCase
	subscriptionController    subscription.Controller
)

func init() {
	db, err := sql.Open("postgres", "postgresql://fain:awurenwae@localhost:5432/banda_lumaksa_development?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	subscriptionRepository = subscription.NewPostgresRepository(db)
	createSubscriptionUseCase = subscription.NewCreateSubscriptionUseCase(subscriptionRepository)
	cancelSubscriptionUseCase = subscription.NewCancelSubscriptionUseCase(subscriptionRepository)
	getSubscriptionUseCase = subscription.NewGetSubscriptionUseCase(subscriptionRepository)
	listSubscriptionsUseCase = subscription.NewListSubscriptionsUseCase(subscriptionRepository)

	subscriptionController = subscription.NewController(
		subscription.With[*subscription.ControllerImpl, subscription.CreateSubscriptionUseCase]("CreateSubscriptionUseCase", createSubscriptionUseCase),
		subscription.With[*subscription.ControllerImpl, subscription.CancelSubscriptionUseCase]("CancelSubscriptionUseCase", cancelSubscriptionUseCase),
		subscription.With[*subscription.ControllerImpl, subscription.GetSubscriptionUseCase]("GetSubscriptionUseCase", getSubscriptionUseCase),
		subscription.With[*subscription.ControllerImpl, subscription.ListSubscriptionsUseCase]("ListSubscriptionsUseCase", listSubscriptionsUseCase),
	)
}

func main() {
	s := server.New(
		server.WithControllers(
			subscriptionController,
		),
	)

	if err := s.Start(); err != nil {
		log.Fatal("Failed to start server.")
	}
}
