package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"golang.org/x/crypto/acme/autocert"

	"github.com/sobat-binangkit/webhook/plugins"
)

func changeFileExtension(filename, newext string) string {
	ext := filepath.Ext(filename)
	name := filename[0 : len(filename)-len(ext)]
	return name + newext
}

func getConfigMap(filename string) (configMap map[string]map[string]string) {
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

						handlerName, ok := config[method]
						if ok {
							handler, ok := handlers[handlerName]
							if !ok {
								handler = getEchoHandlerFunc(p, handlerName)
								handlers[handlerName] = handler
							}
							e.Add(method, path, handler)
						}
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

func getEchoHandlerFunc(p *plugin.Plugin, handlerName string) (wrapper echo.HandlerFunc) {

	wrapper = nil

	sym, err := p.Lookup(handlerName)
	if err == nil {

		ok := true

		wrapper, ok = sym.(func(c echo.Context) error)

		if !ok {
			handler, ok := sym.(func(params map[string]interface{}) (interface{}, int, error))
			if ok {
				pmh := new(plugins.ParamMapHandler)
				pmh.Handler = handler
				wrapper = pmh.Wrapper
			}
		}

		if !ok {
			handler, ok := sym.(func(inp interface{}) (interface{}, int, error))
			if ok {
				sbh := new(plugins.SingleBindingHandler)
				sbh.Handler = handler
				wrapper = sbh.Wrapper
			}
		}

		if !ok {
			handler, ok := sym.(func(params map[string]interface{}, inp interface{}) (interface{}, int, error))
			if ok {
				sbpmh := new(plugins.SingleBindingWithParamMapHandler)
				sbpmh.Handler = handler
				wrapper = sbpmh.Wrapper
			}
		}

		if !ok {
			fmt.Printf("%s not webhook handler function.\n", handlerName)
		}

	} else {

		fmt.Printf("Lookup Error = %s\n", err.Error())

	}

	return wrapper
}

func main() {
	e := echo.New()

	level := os.Getenv("DEBUG_LEVEL")
	if level == "" {
		level = "INFO"
	}

	debugLvl := log.INFO
	switch level {
	case "DEBUG":
		debugLvl = log.DEBUG
	case "INFO":
		debugLvl = log.INFO
	case "WARN":
		debugLvl = log.WARN
	case "ERROR":
		debugLvl = log.ERROR
	case "OFF":
		debugLvl = log.OFF
	default:
		debugLvl = log.INFO
	}

	e.Logger.SetLevel(debugLvl)

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
