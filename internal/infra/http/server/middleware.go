package http_server

import (
	"net/http"

	common_errors "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (server *Server) HTTPErrorHandler(err error, c echo.Context) {
	if val, ok := err.(*common_errors.Error); ok {
		c.JSON(val.Code, echo.Map{
			"error": val,
		})

		return
	}

	code := http.StatusInternalServerError
	if e, ok := err.(*echo.HTTPError); ok {
		code = e.Code
	}

	if code == http.StatusNotFound {
		c.JSON(code, echo.Map{
			"error": common_errors.ErrNotFound.Format(c.Request().Method, c.Request().URL),
		})

		return
	}

	c.JSON(code, echo.Map{
		"error": common_errors.ErrInternalServer,
	})
}

func (server *Server) RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			args := []any{
				logger.String("request-id", v.RequestID),
				logger.String("method", v.Method),
				logger.String("uri", v.URI),
				logger.Int("status", v.Status),
				logger.Int("latency", int(v.Latency)),
			}
			
			if v.Error == nil {
				server.Logger.Info("REQUEST", args...)
			} else {
				if v.Status == http.StatusNotFound {
					server.Logger.Warn("NOT_FOUND", args...)
					return nil
				}
				if val, ok := v.Error.(*common_errors.Error); ok {
					server.Logger.Warn(val.Reason, args...)
					return nil
				}
				args = append(args, logger.String("error", v.Error.Error()))
				server.Logger.Warn("INTERNAL_SERVER_ERROR", args...)
			}
			return nil
		},
		HandleError:      false,
		LogLatency:       true,
		LogProtocol:      false,
		LogRemoteIP:      false,
		LogHost:          false,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       false,
		LogRoutePath:     false,
		LogRequestID:     true,
		LogReferer:       false,
		LogUserAgent:     false,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
	})
}
