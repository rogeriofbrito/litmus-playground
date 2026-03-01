package util

import (
	"context"

	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

type txKey struct{}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func GetTx(ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if !ok {
		return nil, stacktrace.NewError("Failed to get tx from context")
	}
	return tx, nil
}
