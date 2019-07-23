package main

import (
	"fmt"
	"net/http"
)

type (
	HelloConfig struct {
		Greeting string
		Name     string
	}
)

func GetHelloHandler(cfg HelloConfig) (interface{}, int, error) {

	res := make(map[string]string)

	greeting := cfg.Greeting
	if greeting == "" {
		greeting = "Hello"
	}

	name := cfg.Name
	if name == "" {
		name = "World"
	}

	res["result"] = fmt.Sprintf("%s %s !", greeting, name)

	return res, http.StatusOK, nil

}
