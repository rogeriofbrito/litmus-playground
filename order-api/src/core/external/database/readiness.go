package external_database

import "context"

type IReadiness interface {
	OpenConn(ctx context.Context) error
}
