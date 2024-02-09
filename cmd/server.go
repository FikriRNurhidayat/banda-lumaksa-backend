package main

import (
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/config"
	http_server "github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/http/server"
	_ "github.com/lib/pq"
)

func main() {
	config.Init()
	
	srv := http_server.New()

	if err := srv.Start(); err != nil {
		panic(err.Error())
	}
}
