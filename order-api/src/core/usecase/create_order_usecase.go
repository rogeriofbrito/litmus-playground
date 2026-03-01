package usecase

import (
	"context"
	"time"

	"github.com/palantir/stacktrace"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"
	external_database "github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/external/database"
)

func NewCreateOrderUseCase(od external_database.IOrderDatabase) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		od: od,
	}
}

type CreateOrderUseCase struct {
	od external_database.IOrderDatabase
}

func (uc CreateOrderUseCase) Execute(ctx context.Context, order *domain.OrderDomain) (*domain.OrderDomain, error) {
	order.OrderDate = time.Now()

	order, err := uc.od.Save(ctx, order)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed on call database")
	}

	return order, nil
}
