package usecase

import (
	"context"

	"github.com/palantir/stacktrace"
	external_database "github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/external/database"
)

func NewCheckDatabaseUseCase(r external_database.IReadiness) *CheckDatabaseUseCase {
	return &CheckDatabaseUseCase{
		r: r,
	}
}

type CheckDatabaseUseCase struct {
	r external_database.IReadiness
}

func (uc CheckDatabaseUseCase) Execute(ctx context.Context) error {
	if err := uc.r.OpenConn(ctx); err != nil {
		return stacktrace.Propagate(err, "Failed on call database")
	}

	return nil
}
