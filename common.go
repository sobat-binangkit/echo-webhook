package webhook

import (
	"github.com/labstack/echo"
)

type Initiator interface {
	Setup(e *echo.Echo)
}

type HandlerMapper interface {
	GetHandler(path string) echo.HandlerFunc
}
