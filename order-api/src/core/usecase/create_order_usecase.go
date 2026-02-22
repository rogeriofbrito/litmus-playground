package usecase

import (
	"context"
	"time"

	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"
	external_database "github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/external/database"
)

type CreateOrderUseCase struct {
	OrderDatabase external_database.IOrderDatabase
}

func (uc CreateOrderUseCase) Execute(ctx context.Context, order *domain.OrderDomain) (*domain.OrderDomain, error) {
	order.OrderDate = time.Now()

	order, err := uc.OrderDatabase.Save(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
