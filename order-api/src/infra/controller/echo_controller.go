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

func (ctl EchoController) Liveness(_ context.Context) *HealthResponseModel {
	log.Info("Handling liveness request")

	return &HealthResponseModel{
		Status: "UP",
	}
}

func (ctl EchoController) Readiness(ctx context.Context) (*HealthResponseModel, error) {
	log.Info("Handling readiness request")

	if err := ctl.CheckDatabaseUseCase.Execute(ctx); err != nil {
		return nil, err
	}

	return &HealthResponseModel{
		Status: "UP",
	}, nil
}

func (ctl EchoController) CreateOrder(_ context.Context, req *CreateOrderRequestModel) (*CreateOrderResponseModel, error) {
	order := domain.OrderDomain{
		CustomerName: req.CustomerName,
	}

	order, err := ctl.CreateOrderUseCase.Execute(order)
	if err != nil {
		return nil, err
	}

	return &CreateOrderResponseModel{
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		OrderDate:    order.OrderDate,
	}, nil
}

func (ctl EchoController) AddItem(_ context.Context, orderID int64, req *AddItemRequestModel) (*AddItemResponseModel, error) {
	item := domain.ItemDomain{
		OrderID:  orderID,
		ItemName: req.ItemName,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	item, err := ctl.AddItemUseCase.Execute(item)
	if err != nil {
		return nil, err
	}

	return &AddItemResponseModel{
		ItemID:   item.ItemID,
		ItemName: item.ItemName,
		Quantity: item.Quantity,
		Price:    item.Price,
	}, nil
}

func (ctl EchoController) Start(_ context.Context) error {
	root := ctl.Echo.Group("/v1")
	order := root.Group("/order")
	item := order.Group("/:orderID/item")
	liveness := root.Group("/liveness")
	readiness := root.Group("/readiness")

	order.POST("", ctl.orderPost)
	item.PUT("", ctl.itemPut)
	liveness.GET("", ctl.livenessGet)
	readiness.GET("", ctl.readinessGet)

	return ctl.Echo.Start(ctl.Port)
}

func (ctl EchoController) orderPost(c echo.Context) error {
	req := &CreateOrderRequestModel{}
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := ctl.Validate.Struct(req); err != nil {
		return err
	}

	res, err := ctl.CreateOrder(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (ctl EchoController) itemPut(c echo.Context) error {
	orderID, err := strconv.ParseInt(c.Param("orderID"), 10, 64)
	if err != nil {
		return err
	}

	req := &AddItemRequestModel{}
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := ctl.Validate.Struct(req); err != nil {
		return err
	}

	res, err := ctl.AddItem(c.Request().Context(), orderID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (ctl EchoController) livenessGet(c echo.Context) error {
	return c.JSON(http.StatusOK, ctl.Liveness(c.Request().Context()))
}

func (ctl EchoController) readinessGet(c echo.Context) error {
	hc, err := ctl.Readiness(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HealthResponseModel{Status: "DOWN"})
	}
	return c.JSON(http.StatusOK, hc)
}
