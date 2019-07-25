package utils

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func DefaultHandler(c echo.Context) error {

	res := make(map[string]interface{})

	res["time"] = time.Now().Format("2 Jan 2006 15:04:05")
	res["info"] = "No handler found. Put handler libraries in handlers directory and restart the server."

	return c.JSONPretty(http.StatusOK, res, " ")

}
