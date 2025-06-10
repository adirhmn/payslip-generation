package reimbursement

import (
	"context"
	"payslip-generation-system/internal/entity/reimbursement"
	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository_db.go -package=mock -destination=mock/repository_db_mock.go
type dbRepoProvider interface {
	InsertReimbursement(ctx context.Context, rmb reimbursement.Reimbursement) (int, error) 
}

type dbRepo struct {
	db *postgres.Postgres
}

func newDBRepo(
	db *postgres.Postgres,
) dbRepoProvider {
	return &dbRepo{
		db: db,
	}
}

func (r *dbRepo) InsertReimbursement(ctx context.Context, rmb reimbursement.Reimbursement) (int, error) {
	var id int
	err := r.db.DB.QueryRowContext(
		ctx,
		queryInsertReimbursement,
		rmb.UserID,
		rmb.PeriodID,
		rmb.Amount,
		rmb.Description,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
