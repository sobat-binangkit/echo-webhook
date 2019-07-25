package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type (
	HelloConfig struct {
		Greeting string
		Name     string
	}
)

func GetHelloHandler(c echo.Context) (err error) {

	cfg := new(HelloConfig)

	err = c.Bind(cfg)
	if err != nil {
		c.Logger().Warnf("Binding error = %s", err.Error())
	}

	if cfg.Greeting == "" {
		cfg.Greeting = "Hello"
	}

	if cfg.Name == "" {
		cfg.Name = "World"
	}

	return c.JSON(http.StatusOK, cfg)

}

func GetParamInfo(c echo.Context) error {

	params := make(map[string]interface{})

	if len(c.ParamNames()) > 0 {

		pathParams := make(map[string]string)
		for _, name := range c.ParamNames() {
			pathParams[name] = c.Param(name)
		}

		if len(pathParams) > 0 {
			params["pathParams"] = pathParams
		}

	}

	queryParams := make(map[string][]string)
	if len(c.QueryParams()) > 0 {

		for name, value := range c.QueryParams() {
			queryParams[name] = value
		}

	}

	if len(queryParams) > 0 {
		params["queryParams"] = queryParams
	}

	if _, ok := queryParams["prettify"]; ok {
		return c.JSONPretty(http.StatusOK, params, " ")
	} else {
		return c.JSON(http.StatusOK, params)
	}

}

func PostParamInfo(c echo.Context) error {

	params := make(map[string]interface{})

	values, err := c.FormParams()
	if err == nil {

		formParams := make(map[string][]string)
		for name, value := range values {
			formParams[name] = value
		}

		if len(formParams) > 0 {
			params["formParams"] = formParams
		}

	} else {
		c.Logger().Warnf("Form handling error : %s", err.Error())
	}

	if len(c.ParamNames()) > 0 {

		pathParams := make(map[string]string)
		for _, name := range c.ParamNames() {
			pathParams[name] = c.Param(name)
		}

		if len(pathParams) > 0 {
			params["pathParams"] = pathParams
		}

	}

	queryParams := make(map[string][]string)
	if len(c.QueryParams()) > 0 {

		for name, value := range c.QueryParams() {
			queryParams[name] = value
		}

	}

	if len(queryParams) > 0 {
		params["queryParams"] = queryParams
	}

	if _, ok := queryParams["prettify"]; ok {
		return c.JSONPretty(http.StatusOK, params, " ")
	} else {
		return c.JSON(http.StatusOK, params)
	}

}

func main() {
}
