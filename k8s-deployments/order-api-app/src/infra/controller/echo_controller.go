package infra_controller

import (
	"net/http"
	"strconv"

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
	AddItemUsecase     usecase.AddItemUseCase
}

func (c EchoController) Health() HealthResponseModel {
	return HealthResponseModel{
		Status: "UP",
	}
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

	item, err := c.AddItemUsecase.Execute(item)
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
	health := root.Group("/health")

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

	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, controller.Health())
	})

	return controller.Echo.Start(controller.Port)
}
