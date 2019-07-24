package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"golang.org/x/crypto/acme/autocert"

	"github.com/sobat-binangkit/webhook/plugins"
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

	handlerMap := make(map[string]echo.HandlerFunc)
	handlerMap, err := plugins.LoadEchoHandlerFuncs(e, handlerMap, "./handlers")

	if err != nil {
		e.Logger.Warnf("Error in handler loading : %s", err)
	}

	e.Logger.Info("Available path(s) : ")

	routes := e.Routes()
	for i := 0; i < len(routes); i++ {
		e.Logger.Infof("%s[%s] = %s", routes[i].Path, routes[i].Method, routes[i].Name)
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
