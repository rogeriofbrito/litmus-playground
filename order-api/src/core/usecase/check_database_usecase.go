package usecase

import (
	"context"

	external_database "github.com/rogeriofbrito/kubernetes-playground/order-api/src/core/external/database"
)

type CheckDatabaseUseCase struct {
	IReadines external_database.IReadiness
}

func (uc CheckDatabaseUseCase) Execute(ctx context.Context) error {
	return uc.IReadines.OpenConn(ctx)
}
