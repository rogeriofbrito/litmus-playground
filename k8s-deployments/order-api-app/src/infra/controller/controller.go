package infra_controller

type IController interface {
	CreateOrder(body CreateOrderRequestModel) (CreateOrderResponseModel, error)
	Start() error
}
