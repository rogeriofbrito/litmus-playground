package usecase

import (
	"context"

	"github.com/palantir/stacktrace"
	external_database "github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/external/database"
)

type CheckDatabaseUseCase struct {
	IReadines external_database.IReadiness
}

func (uc CheckDatabaseUseCase) Execute(ctx context.Context) error {
	if err := uc.IReadines.OpenConn(ctx); err != nil {
		return stacktrace.Propagate(err, "Failed on call database")
	}

	return nil
}
