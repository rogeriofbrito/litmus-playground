package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/usecase"
	infra_controller "github.com/rogeriofbrito/kubernetes-playground/order-api/src/infra/controller"
	infra_database "github.com/rogeriofbrito/kubernetes-playground/order-api/src/infra/database"
	log "github.com/sirupsen/logrus"
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
		Port:               fmt.Sprintf(":%s", os.Getenv("PORT")),
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
			log.Errorf("Request timeout: %s", c.Path())
		},
		Timeout: 30 * time.Second,
	}))
	return e
}
