package external_database

import "github.com/rogeriofbrito/litmus-playground/order-api/src/core/domain"

type IItemDatabase interface {
	Save(item domain.ItemDomain) (domain.ItemDomain, error)
}
