package infra_database

import (
	"context"

	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

func NewPostgresReadiness(db *gorm.DB) *PostgresReadiness {
	return &PostgresReadiness{
		db: db,
	}
}

type PostgresReadiness struct {
	db *gorm.DB
}

func (r PostgresReadiness) OpenConn(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return stacktrace.Propagate(err, "Failed to get sql DB")
	}

	if err := sqlDB.Ping(); err != nil {
		return stacktrace.Propagate(err, "Failed to ping database")
	}

	return nil
}
