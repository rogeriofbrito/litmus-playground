package main

import (
	"log"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
		Echo:               newEchoClient(),
		Port:               os.Getenv("PORT"),
		CreateOrderUseCase: co,
		AddItemUsecase:     ai,
	}
	if err := controller.Start(); err != nil {
		panic(err)
	}
}

func newEchoClient() *echo.Echo {
	e := echo.New()
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "custom timeout error message returns to client",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Println(c.Path())
		},
		Timeout: 30 * time.Second,
	}))
	return e
}
