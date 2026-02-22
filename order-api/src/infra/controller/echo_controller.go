package infra_controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/usecase"
	log "github.com/sirupsen/logrus"
)

type EchoController struct {
	Validate             *validator.Validate
	Echo                 *echo.Echo
	Port                 string
	CreateOrderUseCase   usecase.CreateOrderUseCase
	AddItemUseCase       usecase.AddItemUseCase
	CheckDatabaseUseCase usecase.CheckDatabaseUseCase
}

func (c EchoController) Liveness(_ context.Context) HealthResponseModel {
	log.Info("Handling liveness request")
	return HealthResponseModel{
		Status: "UP",
	}
}

func (c EchoController) Readiness(ctx context.Context) (HealthResponseModel, error) {
	log.Info("Handling readiness request")

	if err := c.CheckDatabaseUseCase.Execute(ctx); err != nil {
		return HealthResponseModel{}, err
	}

	return HealthResponseModel{
		Status: "UP",
	}, nil
}

func (c EchoController) CreateOrder(req CreateOrderRequestModel) (CreateOrderResponseModel, error) {
	order := domain.OrderDomain{
		CustomerName: req.CustomerName,
	}

	order, err := c.CreateOrderUseCase.Execute(order)
	if err != nil {
		return CreateOrderResponseModel{}, err
	}

	return CreateOrderResponseModel{
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		OrderDate:    order.OrderDate,
	}, nil
}

func (c EchoController) AddItem(orderID int64, req AddItemRequestModel) (AddItemResponseModel, error) {
	item := domain.ItemDomain{
		OrderID:  orderID,
		ItemName: req.ItemName,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	item, err := c.AddItemUseCase.Execute(item)
	if err != nil {
		return AddItemResponseModel{}, err
	}

	return AddItemResponseModel{
		ItemID:   item.ItemID,
		ItemName: item.ItemName,
		Quantity: item.Quantity,
		Price:    item.Price,
	}, nil
}

func (controller EchoController) Start() error {
	root := controller.Echo.Group("/v1")
	order := root.Group("/order")
	item := order.Group("/:orderID/item")
	liveness := root.Group("/liveness")
	readiness := root.Group("/readiness")

	order.POST("", func(c echo.Context) error {
		// bind request body
		req := CreateOrderRequestModel{}
		if err := c.Bind(&req); err != nil {
			return err
		}

		// validate request body
		if err := controller.Validate.Struct(req); err != nil {
			return err
		}

		// evaluate
		res, err := controller.CreateOrder(req)
		if err != nil {
			return err
		}

		// response
		return c.JSON(http.StatusOK, res)
	})

	item.PUT("", func(c echo.Context) error {
		// bind params
		orderID, err := strconv.ParseInt(c.Param("orderID"), 10, 64)
		if err != nil {
			return err
		}

		// bind request body
		req := AddItemRequestModel{}
		if err := c.Bind(&req); err != nil {
			return err
		}

		// validate request body
		if err := controller.Validate.Struct(req); err != nil {
			return err
		}

		// evaluate
		res, err := controller.AddItem(orderID, req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})

	liveness.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, controller.Liveness(c.Request().Context()))
	})

	readiness.GET("", func(c echo.Context) error {
		hc, err := controller.Readiness(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, HealthResponseModel{Status: "DOWN"})
		}
		return c.JSON(http.StatusOK, hc)
	})

	return controller.Echo.Start(controller.Port)
}
