package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/sobat-binangkit/webhook/handlers"
)

func main() {
	e := echo.New()

	path := os.Getenv("WEBHOOK_PATH")
	if path == "" {
		path = "/"
	}

	fmt.Printf("Setting handler for %s", path)

	e.GET(path, handlers.HelloWorldHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
