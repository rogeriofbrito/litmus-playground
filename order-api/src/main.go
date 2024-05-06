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
	pid := infra_database.PostgresItemDatabase{}

	co := usecase.CreateOrderUseCase{
		OrderDatabase: pod,
	}

	ai := usecase.AddItemUseCase{
		OrderDatabase: pod,
		ItemDatabase:  pid,
	}

	controller := infra_controller.EchoController{
		Validate:           validator.New(),
		Echo:               echo.New(),
		Port:               os.Getenv("PORT"),
		CreateOrderUseCase: co,
		AddItemUsecase:     ai,
	}
	if err := controller.Start(); err != nil {
		panic(err)
	}
}
