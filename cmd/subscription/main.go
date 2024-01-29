package main

import (
	"context"
	"database/sql"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	db, err := sql.Open("postgres", "postgresql://fain:awurenwae@127.0.0.1:5432/banda_lumaksa_development?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	subscriptionCommand := subscription.InitializeCommands(db)
	if err := subscriptionCommand.ChargeSubscriptions(context.TODO()); err != nil {
		log.Fatal(err.Error())
	}
}
