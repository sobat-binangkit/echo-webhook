package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/sobat-binangkit/echo-setup/handlers"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	e := echo.New()

	path := os.Getenv("WEBHOOK_PATH")
	if path == "" {
		path = "/"
	}

	datapath := os.Getenv("DATA_PATH")
	if datapath == "" {
		datapath = "./data"
	}

	domain := os.Getenv("DOMAIN_NAME")

	httpaddr := os.Getenv("HTTP_PORT")
	if httpaddr == "" {
		httpaddr = "8080"
	}
	httpaddr = ":" + httoaddr

	ssladdr := os.Getenv("HTTPS_PORT")
	if ssladdr == "" {
		ssladdr = "8443"
	}
	ssladdr = ":" + ssladdr

	fmt.Printf("Domain name = %s\n", domain)
	fmt.Printf("Data directory = %s\n", datapath)
	fmt.Printf("Setting handler for %s\n", path)

	e.GET(path, handlers.HelloWorldHandler)

	if domain == "" {
		e.Logger.Fatal(e.Start(httpaddr))
	} else {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(domain)
		// Cache certificates
		e.AutoTLSManager.Cache = autocert.DirCache(datapath)
		e.Logger.Fatal(e.StartAutoTLS(ssladdr))

	}

}
