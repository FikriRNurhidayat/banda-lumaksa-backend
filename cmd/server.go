package main

import (
	"fmt"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription"
)

func main() {
	createSubscriptionUseCase := subscription.NewCreateSubscriptionUseCase()

	ctl := subscription.NewController(
		subscription.With[*subscription.ControllerImpl, subscription.CreateSubscriptionUseCase]("CreateSubscriptionUseCase", createSubscriptionUseCase),
	)

	fmt.Println(ctl)
}
