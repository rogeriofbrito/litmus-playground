package external_database

import (
	"context"

	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"
)

type IOrderDatabase interface {
	Save(ctx context.Context, order *domain.OrderDomain) (*domain.OrderDomain, error)
	Count(ctx context.Context, orderID int64) (int64, error)
}
