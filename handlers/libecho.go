package main

import (
	"net/http"

	"github.com/labstack/echo"

	"fmt"
)

func GetEcho(c echo.Context) error {
	name := c.QueryParam("guestName")
	if name == "" {
		name = "World"
	}

	return c.String(http.StatusOK, fmt.Sprintf("Hello %s !", name))
}
