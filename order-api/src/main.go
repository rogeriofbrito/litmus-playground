package main

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	infra_controller "github.com/rogeriofbrito/litmus-playground/order-api/src/infra/controller"
)

func main() {
	controller := infra_controller.EchoController{
		Validate: validator.New(),
		Echo:     echo.New(),
		Port:     "127.0.0.1:8080",
	}
	if err := controller.Start(); err != nil {
		panic(err)
	}
}
