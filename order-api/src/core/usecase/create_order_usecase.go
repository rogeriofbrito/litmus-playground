package usecase

import "github.com/rogeriofbrito/litmus-playground/order-api/src/core/domain"

type CreateOrderUseCase struct{}

func (uc CreateOrderUseCase) Execute(order domain.OrderDomain) (domain.OrderDomain, error) {
	return domain.OrderDomain{}, nil
}
