package handlers

import (
	"github.com/labstack/echo"
)

type (
	WebhookManager struct {
		handlers map[string]echo.HandlerFunc
	}
)

func (wm *WebhookManager) Setup(echo *echo.Echo) {

}

func (wm *WebhookManager) GetHandler(path string) echo.HandlerFunc {
	return wm.handlers[path]
}
