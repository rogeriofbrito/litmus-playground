package infra_controller

import "time"

type CreateOrderRequestModel struct {
	CustomerName string `json:"customerName" validate:"required"`
}

type CreateOrderResponseModel struct {
	OrderID      int64     `json:"orderID"`
	CustomerName string    `json:"customerName"`
	OrderDate    time.Time `json:"orderDate"`
}
