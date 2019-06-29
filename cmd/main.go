package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/sobat-binangkit/webhook/handlers"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	e := echo.New()

	path := os.Getenv("WEBHOOK_PATH")
	if path == "" {
		path = "/"
	}

	domain := os.Getenv("DOMAIN_NAME")

	fmt.Printf("Domain name = %s\n", domain)
	fmt.Printf("Setting handler for %s\n", path)

	e.GET(path, handlers.HelloWorldHandler)

	if domain == "" {
		e.Logger.Fatal(e.Start(":8080"))
	} else {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("webhook.id")
		// Cache certificates
		e.AutoTLSManager.Cache = autocert.DirCache(".")
		e.Logger.Fatal(e.StartAutoTLS(":443"))

	}

}
