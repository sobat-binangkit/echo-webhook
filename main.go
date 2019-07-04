package main

import (
	"fmt"
	"os"
	"plugin"

	"github.com/labstack/echo"

	"golang.org/x/crypto/acme/autocert"
)

func getEchoHandlerFuncs(libname string, handlerNames map[string]string) map[string]echo.HandlerFunc {
	handlers := make(map[string]echo.HandlerFunc)

	p, err := plugin.Open(libname)
	if err == nil {

		for path, handlerName := range handlerNames {
			sym, err := p.Lookup(handlerName)
			if err == nil {
				handler, ok := sym.(func(c echo.Context) error)
				if ok {
					handlers[path] = handler
				} else {
					fmt.Printf("%s not echo.HandlerFunc\n", handlerName)
				}
			} else {
				fmt.Printf("Lookup Error = %s\n", err.Error())
			}
		}

	} else {
		fmt.Printf("Plugin.Open Error = %s\n", err.Error())
	}

	return handlers
}

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
	httpaddr = ":" + httpaddr

	ssladdr := os.Getenv("HTTPS_PORT")
	if ssladdr == "" {
		ssladdr = "8443"
	}
	ssladdr = ":" + ssladdr

	fmt.Printf("Domain name = %s\n", domain)
	fmt.Printf("Data directory = %s\n", datapath)
	fmt.Printf("Setting handler for %s\n", path)

	handlerNames := map[string]string{
		"/": "GetHello",
	}

	fmt.Printf("handlerNames = %-v\n", handlerNames)

	handlers := getEchoHandlerFuncs("./handlers/libhello.so", handlerNames)

	fmt.Printf("handlers = %-v\n", handlers)

	for pathname, handler := range handlers {
		e.GET(pathname, handler)
	}

	if domain == "" {
		e.Logger.Fatal(e.Start(httpaddr))
	} else {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(domain)
		// Cache certificates
		e.AutoTLSManager.Cache = autocert.DirCache(datapath)
		e.Logger.Fatal(e.StartAutoTLS(ssladdr))

	}

}
