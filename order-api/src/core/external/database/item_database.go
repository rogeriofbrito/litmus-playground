package external_database

import (
	"context"

	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"
)

type IItemDatabase interface {
	Save(ctx context.Context, item *domain.ItemDomain) (*domain.ItemDomain, error)
}
