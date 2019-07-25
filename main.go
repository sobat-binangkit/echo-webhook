package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"golang.org/x/crypto/acme/autocert"

	"github.com/sobat-binangkit/webhook/plugins"
	"github.com/sobat-binangkit/webhook/utils"
)

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

	datadir := os.Getenv("DATA_DIR")
	if datadir == "" {
		datadir = "./data"
	}

	plugindir := os.Getenv("PLUGIN_DIR")
	if plugindir == "" {
		plugindir = "./handlers"
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

	e.Logger.Infof("Domain name = %s", domain)
	e.Logger.Infof("Data directory = %s", datadir)

	handlerMap := make(map[string]echo.HandlerFunc)
	handlerMap, err := plugins.LoadEchoHandlerFuncs(e, handlerMap, plugindir)

	if err != nil {
		e.Logger.Warnf("Error in handler loading : %s", err)
	}

	e.Logger.Info("Available path(s) : ")

	routes := e.Routes()
	for i := 0; i < len(routes); i++ {
		if level == "DEBUG" {
			e.Logger.Debugf("%s[%s] = %s", routes[i].Path, routes[i].Method, routes[i].Name)
		} else {
			e.Logger.Infof("%s[%s] = %s", routes[i].Path, routes[i].Method)
		}
	}

	if len(routes) == 0 {
		e.GET("/", utils.DefaultHandler)
	}

	if domain == "" {
		e.Logger.Fatal(e.Start(httpaddr))
	} else {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(domain)
		// Cache certificates
		e.AutoTLSManager.Cache = autocert.DirCache(datadir)
		e.Logger.Fatal(e.StartAutoTLS(ssladdr))

	}

}
