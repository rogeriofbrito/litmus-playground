package infra_controller

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/rogeriofbrito/litmus-playground/order-api/src/core/domain"
	"github.com/rogeriofbrito/litmus-playground/order-api/src/core/usecase"
)

type EchoController struct {
	Validate           *validator.Validate
	Echo               *echo.Echo
	Port               string
	CreateOrderUseCase usecase.CreateOrderUseCase
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

func (controller EchoController) Start() error {
	root := controller.Echo.Group("/v1")
	order := root.Group("/order")

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

	return controller.Echo.Start(controller.Port)
}
