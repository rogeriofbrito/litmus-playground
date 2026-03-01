package infra_database

import (
	"context"

	"github.com/palantir/stacktrace"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"
	infra_error "github.com/rogeriofbrito/kubernetes-playground/order-api/src/infra/error"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/infra/util"
	"gorm.io/gorm"
)

func NewPostgresOrderDatabase(db *gorm.DB) *PostgresOrderDatabase {
	return &PostgresOrderDatabase{
		db: db,
	}
}

type PostgresOrderDatabase struct {
	db *gorm.DB
}

func (d PostgresOrderDatabase) Save(ctx context.Context, order *domain.OrderDomain) (*domain.OrderDomain, error) {
	tx, err := util.GetTx(ctx, d.db)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to get tx from context")
	}

	model := OrderModel{
		OrderID:      order.OrderID,
		CustomerName: order.CustomerName,
		OrderDate:    order.OrderDate,
	}

	result := tx.WithContext(ctx).Table("order").Create(&model)
	err = result.Error
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to execute query")
	}

	if result.RowsAffected == 0 {
		return nil, stacktrace.NewErrorWithCode(infra_error.EcodeQueryNotReturnValues, "Query doesn't return values")
	}

	order.OrderID = model.OrderID
	order.CustomerName = model.CustomerName
	order.OrderDate = model.OrderDate

	return order, nil
}

func (d PostgresOrderDatabase) Count(ctx context.Context, orderID int64) (int64, error) {
	tx, err := util.GetTx(ctx, d.db)
	if err != nil {
		return 0, stacktrace.Propagate(err, "Failed to get tx from context")
	}

	var count int64
	result := tx.WithContext(ctx).Table("order").Where("order_id = ?", orderID).Count(&count)
	err = result.Error
	if err != nil {
		return 0, stacktrace.Propagate(err, "Failed to execute query")
	}

	return count, nil
}
