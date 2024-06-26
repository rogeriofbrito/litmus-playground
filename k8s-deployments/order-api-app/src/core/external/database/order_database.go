package external_database

import "github.com/rogeriofbrito/litmus-playground/order-api/src/core/domain"

type IOrderDatabase interface {
	Save(order domain.OrderDomain) (domain.OrderDomain, error)
	Count(orderID int64) (int64, error)
}
