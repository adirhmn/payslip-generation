package ping

import (
	"context"

	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository_db.go -package=mock -destination=mock/repository_db_mock.go
type dbRepoProvider interface {
	Ping(ctx context.Context) error
}

type dbRepo struct {
	db    *postgres.Postgres
}

func newDBRepo(
	db *postgres.Postgres,
) dbRepoProvider {
	return &dbRepo{
		db: db,
	}
}

func (r *dbRepo) Ping(ctx context.Context) error {
	err := r.db.DB.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}