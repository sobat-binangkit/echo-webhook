package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/labstack/echo"
)

func changeFileExtension(filename, newext string) string {
	ext := filepath.Ext(filename)
	name := filename[0 : len(filename)-len(ext)]
	return name + newext
}

func getConfigMap(filename string) (configMap map[string]map[string]string, err error) {
	configMap = make(map[string]map[string]string)

	file, err := os.Open(filename)
	if err == nil {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&configMap)
	}

	return configMap, err
}

func getEchoHandlerFunc(p *plugin.Plugin, handlerName string) (wrapper echo.HandlerFunc, err error) {

	sym, err := p.Lookup(handlerName)

	if err == nil {

		ok := true

		wrapper, ok = sym.(func(echo.Context) error)

		if !ok {
			err = fmt.Errorf("%s not webhook handler function.\n", handlerName)
		}

	} else {

		wrapper = nil

	}

	return wrapper, err
}

func LoadEchoHandlerFuncs(e *echo.Echo, handlers map[string]echo.HandlerFunc, dirname string) (map[string]echo.HandlerFunc, error) {

	methods := []string{"POST", "GET", "PUT", "PATCH", "DELETE"}

	filenames, err := filepath.Glob(dirname + "/*.json")

	if err == nil {

		e.Logger.Infof("Loading = %-v ...\n", filenames)

		for _, filename := range filenames {

			libname := changeFileExtension(filename, ".so")
			e.Logger.Debugf("libname = %s\n", libname)

			p, err := plugin.Open(libname)

			if err == nil {

				configMap, _ := getConfigMap(filename)
				e.Logger.Debugf("configs : %-v\n", configMap)

				for path, config := range configMap {

					for _, method := range methods {

						handlerName, ok := config[method]
						if ok {
							handler, ok := handlers[handlerName]

							if !ok {
								handler, err = getEchoHandlerFunc(p, handlerName)

								if err == nil {
									handlers[handlerName] = handler
								}
							}

							if handler != nil {
								e.Add(method, path, handler)
							}
						}
					}

				}

			} else {

				e.Logger.Debugf("Fail to open %s [%s]\n", libname, err.Error())

			}

		}

	}

	return handlers, err

}
