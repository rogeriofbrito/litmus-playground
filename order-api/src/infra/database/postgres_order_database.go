package infra_database

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/rogeriofbrito/litmus-playground/order-api/src/core/domain"
)

type PostgresOrderDatabase struct{}

func (d PostgresOrderDatabase) Save(order domain.OrderDomain) (domain.OrderDomain, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return domain.OrderDomain{}, err
	}
	defer conn.Close(context.Background())

	insert := `
	INSERT INTO "order" (
	  	  customer_name
		, order_date
	) VALUES (
		  $1
		, $2
	)
	returning 
		  order_id
		, customer_name
		, order_date`

	rows, err := conn.Query(context.Background(), insert, order.CustomerName, order.OrderDate)
	if err != nil {
		return domain.OrderDomain{}, err
	}

	if rows.Next() {
		err = rows.Scan(&order.OrderID, &order.CustomerName, &order.OrderDate)
		if err != nil {
			return domain.OrderDomain{}, err
		}
	} else {
		return domain.OrderDomain{}, errors.New("PostgresOrderDatabase.Save: insert does't return values")
	}

	return order, nil
}
