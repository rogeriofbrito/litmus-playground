package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/usecase"
	infra_controller "github.com/rogeriofbrito/kubernetes-playground/order-api/src/infra/controller"
	infra_database "github.com/rogeriofbrito/kubernetes-playground/order-api/src/infra/database"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	db, err := gorm.Open(postgres.Open(getConnString()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	configureDBPool(sqlDB)

	pod := infra_database.NewPostgresOrderDatabase(db)
	pid := infra_database.NewPostgresItemDatabase(db)
	pr := infra_database.NewPostgresReadiness(db)
	co := usecase.NewCreateOrderUseCase(pod)
	ai := usecase.NewAddItemUseCase(pod, pid)
	cd := usecase.NewCheckDatabaseUseCase(pr)
	var controller infra_controller.IController = infra_controller.NewEchoController(
		validator.New(),
		newEchoClient(),
		fmt.Sprintf(":%s", os.Getenv("PORT")),
		co,
		ai,
		cd,
		db,
	)

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

	e.Use(echoprometheus.NewMiddleware("echo"))

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		log.Errorf("Failed to handle request at path %s: %s", c.Path(), err.Error())

		if !c.Response().Committed {
			c.JSON(code, map[string]interface{}{
				"message": http.StatusText(code),
			})
		}
	}

	return e
}

func getConnString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
}

func configureDBPool(sqlDB *sql.DB) {
	sqlDB.SetMaxOpenConns(30)                  //TODO: create env var
	sqlDB.SetMaxIdleConns(15)                  //TODO: create env var
	sqlDB.SetConnMaxLifetime(1 * time.Hour)    //TODO: create env var
	sqlDB.SetConnMaxIdleTime(15 * time.Minute) //TODO: create env var
}
