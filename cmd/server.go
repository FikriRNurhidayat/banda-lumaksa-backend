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
	createSubscriptionService subscription.CreateSubscriptionService
	cancelSubscriptionService subscription.CancelSubscriptionService
	getSubscriptionService    subscription.GetSubscriptionService
	listSubscriptionsService  subscription.ListSubscriptionsService
	subscriptionController    subscription.Controller
)

func init() {
	db, err := sql.Open("postgres", "postgresql://fain:awurenwae@127.0.0.1:5432/banda_lumaksa_development?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	subscriptionRepository = subscription.NewPostgresRepository(db)
	createSubscriptionService = subscription.NewCreateSubscriptionService(subscriptionRepository)
	cancelSubscriptionService = subscription.NewCancelSubscriptionService(subscriptionRepository)
	getSubscriptionService = subscription.NewGetSubscriptionService(subscriptionRepository)
	listSubscriptionsService = subscription.NewListSubscriptionsService(subscriptionRepository)

	subscriptionController = subscription.NewController(
		listSubscriptionsService,
		getSubscriptionService,
		createSubscriptionService,
		cancelSubscriptionService,
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
