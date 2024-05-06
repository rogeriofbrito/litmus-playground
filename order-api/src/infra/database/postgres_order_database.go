package infra_database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/rogeriofbrito/litmus-playground/order-api/src/core/domain"
	infra_error "github.com/rogeriofbrito/litmus-playground/order-api/src/infra/error"
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
	) returning 
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
		return domain.OrderDomain{}, infra_error.ErrQueryNotReturnValues
	}

	return order, nil
}

func (d PostgresOrderDatabase) Count(orderID int64) (int64, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return 0, err
	}
	defer conn.Close(context.Background())

	count := `
	SELECT 
		COUNT(order_id) 
	FROM "order" 
	WHERE order_id = $1`

	rows, err := conn.Query(context.Background(), count, orderID)
	if err != nil {
		return 0, err
	}

	var countResult int64
	if rows.Next() {
		err = rows.Scan(&countResult)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, infra_error.ErrQueryNotReturnValues
	}

	return countResult, nil
}
