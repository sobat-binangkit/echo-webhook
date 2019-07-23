package plugins

import (
	"reflect"

	"github.com/labstack/echo"
)

type (
	SignatureType int

	GenericHandler struct {
		Signature SignatureType
	}

	ParamMapHandler struct {
		GenericHandler
		Handler func(params map[string]interface{}) (interface{}, int, error)
	}

	SingleBindingHandler struct {
		GenericHandler
		Handler func(inp interface{}) (interface{}, int, error)
	}

	SingleBindingWithParamMapHandler struct {
		GenericHandler
		Handler func(params map[string]interface{}, inp interface{}) (interface{}, int, error)
	}
)

const (
	EchoContext = iota
	ParamMap
	SingleBinding
	SingleBindingWithParamMap
)

func (s SignatureType) String() string {
	return [...]string{"EchoContext", "ParamMap", "SingleBinding", "SingleBindingWithMap"}[s]
}

func createParamMap(c echo.Context) map[string]interface{} {

	params := make(map[string]interface{})

	pathParams := make(map[string]string)
	for _, name := range c.ParamNames() {
		pathParams[name] = c.Param(name)
	}

	if len(pathParams) > 0 {
		params["pathParams"] = pathParams
	}

	values, err := c.FormParams()
	if err == nil {
		formParams := make(map[string][]string)
		for name, value := range values {
			formParams[name] = value
		}

		if len(formParams) > 0 {
			params["formParams"] = formParams
		}
	}

	queryParams := make(map[string][]string)
	for name, value := range c.QueryParams() {
		queryParams[name] = value
	}

	if len(queryParams) > 0 {
		params["queryParams"] = queryParams
	}

	return params

}

func (h ParamMapHandler) Wrapper(c echo.Context) error {

	params := createParamMap(c)

	res, code, err := h.Handler(params)

	if err != nil {
		return err
	} else {
		return c.JSON(code, res)
	}

}

func (h SingleBindingHandler) Wrapper(c echo.Context) error {

	ht := reflect.TypeOf(h.Handler)
	vt := ht.In(0)
	value := reflect.Zero(vt)

	err := c.Bind(value)

	if err == nil {
		res, code, err := h.Handler(value)
		if err == nil {
			return c.JSON(code, res)
		}
	}

	return err

}

func (h SingleBindingWithParamMapHandler) Wrapper(c echo.Context) error {

	ht := reflect.TypeOf(h.Handler)
	vt := ht.In(1)
	value := reflect.Zero(vt)

	err := c.Bind(value)

	if err == nil {
		params := createParamMap(c)

		res, code, err := h.Handler(params, value)
		if err == nil {
			return c.JSON(code, res)
		}
	}

	return err
}
