package http_server

import (
	common_module "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/module"
	subscription_module "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/module"
	transaction_module "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/transaction/module"
)

var HTTPModules = common_module.HTTPModules{
	subscription_module.HTTP(),
	transaction_module.HTTP(),
}
