package common_controller

import "github.com/labstack/echo/v4"

type Controller interface {
	Register(*echo.Echo)
}
