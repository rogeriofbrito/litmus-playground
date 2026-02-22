package main

import (
	"context"
	"fmt"
	"net/http"
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
	ctx := context.Background()

	pod := infra_database.PostgresOrderDatabase{}
	pid := infra_database.PostgresItemDatabase{}
	pr := infra_database.PostgresReadiness{}

	co := usecase.CreateOrderUseCase{
		OrderDatabase: pod,
	}

	ai := usecase.AddItemUseCase{
		OrderDatabase: pod,
		ItemDatabase:  pid,
	}

	cd := usecase.CheckDatabaseUseCase{
		IReadines: pr,
	}

	var controller infra_controller.IController = infra_controller.EchoController{
		Validate:             validator.New(),
		Echo:                 newEchoClient(),
		Port:                 fmt.Sprintf(":%s", os.Getenv("PORT")),
		CreateOrderUseCase:   co,
		AddItemUseCase:       ai,
		CheckDatabaseUseCase: cd,
	}
	if err := controller.Start(ctx); err != nil {
		log.Fatal(err)
	}
}

func newEchoClient() *echo.Echo {
	e := echo.New()

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request timeout",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Errorf("Timeout for request path %s", c.Path())
		},
		Timeout: 30 * time.Second,
	}))

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		log.Errorf("Error handled for request path %s: %s", c.Path(), err.Error())

		if !c.Response().Committed {
			c.JSON(code, map[string]interface{}{
				"message": http.StatusText(code),
			})
		}
	}

	return e
}
