package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Controller interface {
	Register(*echo.Echo)
}

type Server struct {
	echo *echo.Echo
	Port uint
}

type ServerBuilder func(*Server)

func (s *Server) Start() error {
	return s.echo.Start(fmt.Sprintf(":%d", s.Port))
}

func WithControllers(controllers ...Controller) ServerBuilder {
	return func(s *Server) {
		for _, ctl := range controllers {
			ctl.Register(s.echo)
		}
	}
}

func New(builders ...ServerBuilder) *Server {
	server := &Server{
		echo: echo.New(),
		Port: 3000,
	}

	server.echo.Use(middleware.Logger())
	server.echo.Use(middleware.Recover())

	for _, build := range builders {
		build(server)
	}

	return server
}
