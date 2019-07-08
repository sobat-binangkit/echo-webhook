package main

import (
	"time"

	"net/http"

	"github.com/labstack/echo"

	"fmt"
)

func GetHello(c echo.Context) error {
	name := c.QueryParam("guestName")
	if name == "" {
		name = "World"
	}

	return c.String(http.StatusOK, fmt.Sprintf("Hello %s !\n", name))
}

func GetTimestamp(c echo.Context) error {
	t := time.Now()

	return c.String(http.StatusOK, fmt.Sprintf("Current Timestamp %s\n", t.Format(time.RFC1123Z)))
}
