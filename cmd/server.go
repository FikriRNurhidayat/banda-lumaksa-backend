package main

import (
	"log"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/server"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/transaction"
)

func main() {
	db, err := sql.Open("postgres", "postgresql://fain:awurenwae@127.0.0.1:5432/banda_lumaksa_development?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	subscriptionController := subscription.InitializeHandlers(db)
	transactionController := transaction.InitializeHandlers(db)

	srv := server.New(
		server.WithControllers(
			subscriptionController,
			transactionController,
		),
	)

	if err := srv.Start(); err != nil {
		log.Fatal("Failed to start server.")
	}
}
