package usecase

import (
	"context"

	"github.com/palantir/stacktrace"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"
	core_error "github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/error"
	external_database "github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/external/database"
)

type AddItemUseCase struct {
	OrderDatabase external_database.IOrderDatabase
	ItemDatabase  external_database.IItemDatabase
}

func (uc AddItemUseCase) Execute(ctx context.Context, item *domain.ItemDomain) (*domain.ItemDomain, error) {
	countOrder, err := uc.OrderDatabase.Count(ctx, item.OrderID)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed on call database")
	}

	if countOrder == 0 {
		return nil, stacktrace.NewErrorWithCode(core_error.EcodeOrderNotFound, "Order not found")
	}

	item, err = uc.ItemDatabase.Save(ctx, item)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed on call database")
	}

	return item, nil
}
