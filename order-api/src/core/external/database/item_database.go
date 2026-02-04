package external_database

import "github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"

type IItemDatabase interface {
	Save(item domain.ItemDomain) (domain.ItemDomain, error)
}
