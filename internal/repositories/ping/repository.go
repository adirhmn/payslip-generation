package ping

import (
	"context"

	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository.go -package=mock -destination=mock/repository_mock.go
type PingRepositoryProvider interface {
	Ping(ctx context.Context) error
}

type pingRepository struct {
	db    dbRepoProvider
}

func NewPingRepository(
	db *postgres.Postgres,
) PingRepositoryProvider {
	return &pingRepository{
		db: newDBRepo(
			db,
		),
	}
}

func (r *pingRepository) Ping(ctx context.Context) error {
	// ping database
	err := r.db.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}