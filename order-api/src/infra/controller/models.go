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

type AddItemRequestModel struct {
	ItemName string  `json:"itemName" validate:"required"`
	Quantity int64   `json:"quantity" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
}

type AddItemResponseModel struct {
	ItemID   int64   `json:"itemID"`
	ItemName string  `json:"itemName"`
	Quantity int64   `json:"quantity"`
	Price    float64 `json:"price"`
}
