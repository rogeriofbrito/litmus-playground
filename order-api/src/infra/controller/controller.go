package infra_controller

import "context"

type IController interface {
	Liveness(ctx context.Context) *HealthResponseModel
	Readiness(ctx context.Context) (*HealthResponseModel, error)
	CreateOrder(ctx context.Context, req *CreateOrderRequestModel) (*CreateOrderResponseModel, error)
	AddItem(ctx context.Context, orderID int64, req *AddItemRequestModel) (*AddItemResponseModel, error)
	Start(ctx context.Context) error
}
