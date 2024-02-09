package http_server

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"

	common_module "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/module"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/db"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/logger"
)

type Server struct {
	Echo       *echo.Echo
	Port       uint
	DB         *sql.DB
	Logger     logger.Logger
	Dependency *common_module.RootDependency
}

func (s *Server) Start() error {
	s.Logger.Info("server/START", logger.String("scheme", "http"), logger.String("host", "localhost"), logger.Int("port", int(s.Port)))
	return s.Echo.Start(fmt.Sprintf(":%d", s.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("server/SHUTDOWN", logger.String("scheme", "http"), logger.String("host", "localhost"), logger.Int("port", int(s.Port)))
	return s.Echo.Shutdown(ctx)
}

func New() *Server {
	server := &Server{
		Port:   viper.GetUint("server.port"),
		Echo:   echo.New(),
		DB:     db.New(),
		Logger: logger.New(),
	}

	server.Dependency = common_module.New(server.DB, server.Logger)
	server.Echo.Logger.SetOutput(io.Discard)
	server.Echo.Logger.SetLevel(log.OFF)
	server.Echo.HideBanner = true
	server.Echo.HidePort = true
	server.Echo.DisableHTTP2 = true
	server.Echo.Use(middleware.Secure())
	server.Echo.Use(middleware.Timeout())
	server.Echo.Use(middleware.RequestID())
	server.Echo.Use(server.RequestLogger())
	server.Echo.Use(middleware.Recover())
	server.Echo.GET("/health", server.HealthCheck)
	server.Echo.HTTPErrorHandler = server.HTTPErrorHandler

	for _, mod := range HTTPModules {
		mod.Wire(server.Dependency)
		mod.WireController(server.Echo)
	}

	return server
}
