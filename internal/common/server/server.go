package server

import (
	"fmt"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
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

	server.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		if val, ok := err.(*errors.Error); ok {
			c.JSON(val.Code, echo.Map{
				"error": val,
			})

			return
		}

		server.echo.DefaultHTTPErrorHandler(err, c)
	}

	for _, build := range builders {
		build(server)
	}

	return server
}
