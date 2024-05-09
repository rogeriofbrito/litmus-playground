package usecase

import (
	"github.com/rogeriofbrito/litmus-playground/order-api/src/core/domain"
	core_error "github.com/rogeriofbrito/litmus-playground/order-api/src/core/error"
	external_database "github.com/rogeriofbrito/litmus-playground/order-api/src/core/external/database"
)

type AddItemUseCase struct {
	OrderDatabase external_database.IOrderDatabase
	ItemDatabase  external_database.IItemDatabase
}

func (uc AddItemUseCase) Execute(item domain.ItemDomain) (domain.ItemDomain, error) {
	countOrder, err := uc.OrderDatabase.Count(item.OrderID)
	if err != nil {
		return domain.ItemDomain{}, err
	}

	if countOrder == 0 {
		return domain.ItemDomain{}, core_error.ErrOrderNotFound
	}

	item, err = uc.ItemDatabase.Save(item)
	if err != nil {
		return domain.ItemDomain{}, err
	}

	return item, nil
}
