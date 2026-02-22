package infra_database

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type PostgresReadiness struct{}

func (r PostgresReadiness) OpenConn(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, getConnString())
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	return nil
}
