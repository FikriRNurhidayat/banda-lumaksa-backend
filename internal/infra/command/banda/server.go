package banda_command

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/config"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/config/version"
	http_server "github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/http/server"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/logger"
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run banda server.",
	Long:  `Run banda server.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		log := logger.New(version.Version, version.Build)

		srv, err := http_server.New(log)
		if err != nil {
			log.Fatal("http/FAILURE", logger.String("error", err.Error()))
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		go func() {
			if err := srv.Start(); err != nil && err != http.ErrServerClosed {
				os.Exit(0)
			}
		}()

		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			srv.Logger.Error("http/SHUTDOWN_FAILURE", logger.String("error", err.Error()))
		}
	},
}

