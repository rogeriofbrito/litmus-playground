package main

import (
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/rogeriofbrito/litmus-playground/order-api/src/core/usecase"
	infra_controller "github.com/rogeriofbrito/litmus-playground/order-api/src/infra/controller"
	infra_database "github.com/rogeriofbrito/litmus-playground/order-api/src/infra/database"
)

func main() {
	pod := infra_database.PostgresOrderDatabase{}

	co := usecase.CreateOrderUseCase{
		OrderDatabase: pod,
	}

	controller := infra_controller.EchoController{
		Validate:           validator.New(),
		Echo:               echo.New(),
		Port:               os.Getenv("PORT"),
		CreateOrderUseCase: co,
	}
	if err := controller.Start(); err != nil {
		panic(err)
	}
}
