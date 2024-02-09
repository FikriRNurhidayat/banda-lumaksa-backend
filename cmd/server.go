package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/config"
	http_server "github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/http/server"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/logger"
	_ "github.com/lib/pq"
)

var (
	version string = "dev"
	build   string = fmt.Sprintf("%d", time.Now().UnixMilli())
)

func main() {
	config.Init()

	srv := http_server.New(build, version)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			os.Exit(0)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		srv.Logger.Error("SERVER_SHUTDOWN_FAILURE", logger.String("error", err.Error()))
	}
}
