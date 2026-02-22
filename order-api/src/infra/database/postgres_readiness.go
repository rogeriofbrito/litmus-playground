package infra_database

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/palantir/stacktrace"
)

type PostgresReadiness struct{}

func (r PostgresReadiness) OpenConn(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, getConnString())
	if err != nil {
		return stacktrace.Propagate(err, "Failed to open Postgres connection")
	}
	defer conn.Close(ctx)

	return nil
}
