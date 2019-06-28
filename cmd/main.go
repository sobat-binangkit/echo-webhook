package main

import (
	"github.com/labstack/echo"

	"github.com/sobat-binangkit/webhook/handlers"
)

func main() {
	e := echo.New()

	e.GET("/", handlers.HelloWorldHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
