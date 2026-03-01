package infra_database

import (
	"context"

	"github.com/palantir/stacktrace"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/domain"
	infra_error "github.com/rogeriofbrito/kubernetes-playground/order-api/src/infra/error"
	"github.com/rogeriofbrito/kubernetes-playground/order-api/src/infra/util"
	"gorm.io/gorm"
)

func NewPostgresItemDatabase(db *gorm.DB) *PostgresItemDatabase {
	return &PostgresItemDatabase{
		db: db,
	}
}

type PostgresItemDatabase struct {
	db *gorm.DB
}

func (d PostgresItemDatabase) Save(ctx context.Context, item *domain.ItemDomain) (*domain.ItemDomain, error) {
	tx, err := util.GetTx(ctx, d.db)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to get tx from context")
	}

	model := ItemModel{
		OrderID:  item.OrderID,
		ItemName: item.ItemName,
		Quantity: item.Quantity,
		Price:    item.Price,
	}

	result := tx.WithContext(ctx).Table("item").Create(&model)
	err = result.Error
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to execute query")
	}

	if result.RowsAffected == 0 {
		return nil, stacktrace.NewErrorWithCode(infra_error.EcodeQueryNotReturnValues, "Query doesn't return values")
	}

	item.ItemID = model.ItemID
	item.OrderID = model.OrderID
	item.ItemName = model.ItemName
	item.Quantity = model.Quantity
	item.Price = model.Price

	return item, nil
}
