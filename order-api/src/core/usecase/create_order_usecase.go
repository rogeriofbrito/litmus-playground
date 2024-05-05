package usecase

import (
	"time"

	"github.com/rogeriofbrito/litmus-playground/order-api/src/core/domain"
	external_database "github.com/rogeriofbrito/litmus-playground/order-api/src/core/external/database"
)

type CreateOrderUseCase struct {
	OrderDatabase external_database.IOrderDatabase
}

func (uc CreateOrderUseCase) Execute(order domain.OrderDomain) (domain.OrderDomain, error) {
	order.OrderDate = time.Now()
	order, err := uc.OrderDatabase.Save(order)
	if err != nil {
		return domain.OrderDomain{}, err
	}
	return order, nil
}
