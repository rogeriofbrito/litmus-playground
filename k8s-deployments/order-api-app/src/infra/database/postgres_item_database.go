package infra_database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/rogeriofbrito/litmus-playground/order-api/src/core/domain"
	infra_error "github.com/rogeriofbrito/litmus-playground/order-api/src/infra/error"
)

type PostgresItemDatabase struct{}

func (d PostgresItemDatabase) Save(item domain.ItemDomain) (domain.ItemDomain, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return domain.ItemDomain{}, err
	}
	defer conn.Close(context.Background())

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

	rows, err := conn.Query(context.Background(), insert, item.OrderID, item.ItemName, item.Quantity, item.Price)
	if err != nil {
		return domain.ItemDomain{}, err
	}

	if rows.Next() {
		err = rows.Scan(&item.ItemID, &item.OrderID, &item.ItemName, &item.Quantity, &item.Price)
		if err != nil {
			return domain.ItemDomain{}, err
		}
	} else {
		return domain.ItemDomain{}, infra_error.ErrQueryNotReturnValues
	}

	return item, nil
}
