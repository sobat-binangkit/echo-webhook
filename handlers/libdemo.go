package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type (
	HelloConfig struct {
		Greeting string
		Name     string
	}
)

func GetHelloHandler(c echo.Context) (err error) {

	cfg := new(HelloConfig)

	err = c.Bind(cfg)
	if err != nil {
		c.Logger().Warnf("Binding error : %s", err.Error())
	}

	if cfg.Greeting == "" {
		cfg.Greeting = "Hello"
	}

	if cfg.Name == "" {
		cfg.Name = "World"
	}

	return c.JSON(http.StatusOK, cfg)

}
