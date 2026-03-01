package infra_controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/palantir/stacktrace"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/usecase"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/infra/util"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewEchoController(
	v *validator.Validate,
	echo *echo.Echo,
	port string,
	co *usecase.CreateOrderUseCase,
	ai *usecase.AddItemUseCase,
	cd *usecase.CheckDatabaseUseCase,
	db *gorm.DB,
) *EchoController {
	return &EchoController{
		v:    v,
		echo: echo,
		port: port,
		co:   co,
		ai:   ai,
		cd:   cd,
		db:   db,
	}
}

type EchoController struct {
	v    *validator.Validate
	echo *echo.Echo
	port string
	co   *usecase.CreateOrderUseCase
	ai   *usecase.AddItemUseCase
	cd   *usecase.CheckDatabaseUseCase
	db   *gorm.DB
}

func (ctl EchoController) Liveness(_ context.Context) *HealthResponseModel {
	log.Info("Handling liveness request")

	return &HealthResponseModel{
		Status: "UP",
	}
}

func (ctl EchoController) Readiness(ctx context.Context) (*HealthResponseModel, error) {
	log.Info("Handling readiness request")

	if err := ctl.cd.Execute(ctx); err != nil {
		return nil, stacktrace.Propagate(err, "Failed to execute CheckDatabaseUseCase")
	}

	return &HealthResponseModel{
		Status: "UP",
	}, nil
}

func (ctl EchoController) CreateOrder(ctx context.Context, req *CreateOrderRequestModel) (*CreateOrderResponseModel, error) {
	order := &domain.OrderDomain{
		CustomerName: req.CustomerName,
	}

	order, err := ctl.co.Execute(ctx, order)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to execute CreateOrderUseCase")
	}

	return &CreateOrderResponseModel{
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		OrderDate:    order.OrderDate,
	}, nil
}

func (ctl EchoController) AddItem(ctx context.Context, orderID int64, req *AddItemRequestModel) (*AddItemResponseModel, error) {
	item := &domain.ItemDomain{
		OrderID:  orderID,
		ItemName: req.ItemName,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	item, err := ctl.ai.Execute(ctx, item)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to execute AddItemUseCase")
	}

	return &AddItemResponseModel{
		ItemID:   item.ItemID,
		ItemName: item.ItemName,
		Quantity: item.Quantity,
		Price:    item.Price,
	}, nil
}

func (ctl EchoController) Start(_ context.Context) error {
	root := ctl.echo.Group("/v1")
	order := root.Group("/order")
	item := order.Group("/:orderID/item")
	liveness := root.Group("/liveness")
	readiness := root.Group("/readiness")
	metrics := root.Group("/metrics")

	order.POST("", ctl.orderPost)
	item.PUT("", ctl.itemPut)
	liveness.GET("", ctl.livenessGet)
	readiness.GET("", ctl.readinessGet)
	metrics.GET("", echoprometheus.NewHandler())

	return ctl.echo.Start(ctl.port)
}

func (ctl EchoController) orderPost(c echo.Context) error {
	ctx := c.Request().Context()

	return ctl.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = util.WithTx(ctx, tx)

		req := &CreateOrderRequestModel{}
		if err := c.Bind(req); err != nil {
			return stacktrace.Propagate(err, "Failed to bind request")
		}

		if err := ctl.v.Struct(req); err != nil {
			return stacktrace.Propagate(err, "Failed to validate request")
		}

		res, err := ctl.CreateOrder(ctx, req)
		if err != nil {
			return stacktrace.Propagate(err, "Failed to handle request")
		}

		return c.JSON(http.StatusOK, res)
	})
}

func (ctl EchoController) itemPut(c echo.Context) error {
	ctx := c.Request().Context()

	return ctl.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = util.WithTx(ctx, tx)

		orderID, err := strconv.ParseInt(c.Param("orderID"), 10, 64)
		if err != nil {
			return stacktrace.Propagate(err, "Failed to parse orderID")
		}

		req := &AddItemRequestModel{}
		if err := c.Bind(req); err != nil {
			return stacktrace.Propagate(err, "Failed to bind request")
		}

		if err := ctl.v.Struct(req); err != nil {
			return stacktrace.Propagate(err, "Failed to validate request")
		}

		res, err := ctl.AddItem(ctx, orderID, req)
		if err != nil {
			return stacktrace.Propagate(err, "Failed to handle request")
		}

		return c.JSON(http.StatusOK, res)
	})
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
