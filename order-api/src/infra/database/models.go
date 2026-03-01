package infra_database

import "time"

type OrderModel struct {
	OrderID      int64     `gorm:"column:order_id;primaryKey;autoIncrement"`
	CustomerName string    `gorm:"column:customer_name"`
	OrderDate    time.Time `gorm:"column:order_date"`
}

type ItemModel struct {
	ItemID   int64   `gorm:"column:item_id;primaryKey;autoIncrement"`
	OrderID  int64   `gorm:"column:order_id"`
	ItemName string  `gorm:"column:item_name"`
	Quantity int64   `gorm:"column:quantity"`
	Price    float64 `gorm:"column:price"`
}
