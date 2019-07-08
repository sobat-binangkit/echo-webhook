package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/labstack/echo"

	"golang.org/x/crypto/acme/autocert"
)

func changeFileExtension(filename, newext string) string {
	ext := filepath.Ext(filename)
	name := filename[0 : len(filename)-len(ext)]
	return name + newext
}

func getConfigMap(filename string) (configMap map[string]map[string]string) {
	fmt.Printf("Loading %s...\n", filename)
	configMap = make(map[string]map[string]string)

	file, err := os.Open(filename)
	if err == nil {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&configMap)
	}

	return configMap
}

func loadEchoHandlerFuncs(e *echo.Echo, dirname string) {

	methods := []string{"POST", "GET", "PUT", "PATCH", "DELETE"}

	filenames, err := filepath.Glob(dirname + "/*.json")

	if err == nil {

		e.Logger.Infof("Loading = %-v ...\n", filenames)

		for _, filename := range filenames {

			libname := changeFileExtension(filename, ".so")
			e.Logger.Debugf("libname = %s\n", libname)

			p, err := plugin.Open(libname)

			if err == nil {
				handlers := make(map[string]echo.HandlerFunc)

				configMap := getConfigMap(filename)
				e.Logger.Debugf("configs : %-v\n", configMap)

				for path, config := range configMap {

					for _, method := range methods {

						handlerName := config[method]
						handler, ok := handlers[handlerName]
						if !ok {
							handler = getEchoHandlerFunc(p, handlerName)
							handlers[handlerName] = handler
						}
						e.Add(method, path, handler)
					}

				}
			} else {
				e.Logger.Debugf("Fail to open %s [%s]\n", libname, err.Error())
			}

		}

	} else {
		e.Logger.Debugf("Directory listing error = %s\n", err.Error())
	}

}

func getEchoHandlerFunc(p *plugin.Plugin, handlerName string) (handler echo.HandlerFunc) {

	sym, err := p.Lookup(handlerName)
	if err == nil {
		ok := true
		handler, ok = sym.(func(c echo.Context) error)
		if !ok {
			fmt.Printf("%s not echo.HandlerFunc\n", handlerName)
		}
	} else {
		fmt.Printf("Lookup Error = %s\n", err.Error())
	}

	return handler
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

	e.Logger.Infof("Domain name = %s\n", domain)
	e.Logger.Infof("Data directory = %s\n", datapath)
	e.Logger.Infof("Setting handler for %s\n", path)

	loadEchoHandlerFuncs(e, "./handlers")

	if domain == "" {
		e.Logger.Fatal(e.Start(httpaddr))
	} else {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(domain)
		// Cache certificates
		e.AutoTLSManager.Cache = autocert.DirCache(datapath)
		e.Logger.Fatal(e.StartAutoTLS(ssladdr))

	}

}
