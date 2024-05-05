package domain

import "time"

type OrderDomain struct {
	OrderID      int64
	CustomerName string
	OrderDate    time.Time
	Items        []ItemDomain
}
