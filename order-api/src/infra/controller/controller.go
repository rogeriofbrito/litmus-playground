package infra_controller

import "context"

type IController interface {
	Liveness(ctx context.Context) HealthResponseModel
	Readiness(ctx context.Context) (HealthResponseModel, error)
	CreateOrder(body CreateOrderRequestModel) (CreateOrderResponseModel, error)
	AddItem(orderID int64, req AddItemRequestModel) (AddItemResponseModel, error)
	Start() error
}
