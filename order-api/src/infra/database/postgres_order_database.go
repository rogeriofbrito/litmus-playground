package infra_database

import "github.com/rogeriofbrito/litmus-playground/order-api/src/core/domain"

type PostgresOrderDatabase struct{}

func (d PostgresOrderDatabase) Save(order domain.OrderDomain) (domain.OrderDomain, error) {
	return domain.OrderDomain{}, nil
}
