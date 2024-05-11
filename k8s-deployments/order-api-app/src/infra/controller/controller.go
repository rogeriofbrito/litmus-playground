package infra_controller

type IController interface {
	Health() HealthResponseModel
	CreateOrder(body CreateOrderRequestModel) (CreateOrderResponseModel, error)
	AddItem(orderID int64, req AddItemRequestModel) (AddItemResponseModel, error)
	Start() error
}
