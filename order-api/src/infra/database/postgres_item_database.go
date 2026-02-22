package infra_database

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"
	infra_error "github.com/rogeriofbrito/kubernetes-playground/order-api/src/infra/error"
)

type PostgresItemDatabase struct{}

func (d PostgresItemDatabase) Save(ctx context.Context, item *domain.ItemDomain) (*domain.ItemDomain, error) {
	conn, err := pgx.Connect(ctx, getConnString())
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	insert := `
	INSERT INTO public.item(
		  order_id
		, item_name
		, quantity
		, price
	) VALUES(
		  $1
		, $2
		, $3
		, $4
	) returning 
		  item_id
		, order_id
		, item_name
		, quantity
		, price`

	rows, err := conn.Query(ctx, insert, item.OrderID, item.ItemName, item.Quantity, item.Price)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		err = rows.Scan(&item.ItemID, &item.OrderID, &item.ItemName, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, infra_error.ErrQueryNotReturnValues
	}

	return item, nil
}
