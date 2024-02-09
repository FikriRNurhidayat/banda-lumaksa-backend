package http_server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) HealthCheck(c echo.Context) error {
	if err := s.DB.Ping(); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
