package reimbursement

import (
	"context"

	"payslip-generation-system/internal/entity/reimbursement"
	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository.go -package=mock -destination=mock/repository_mock.go
type ReimbursementRepositoryProvider interface {
	InsertReimbursement(ctx context.Context, rmb reimbursement.Reimbursement) (int, error)
}

type reimbursementRepository struct {
	db    dbRepoProvider
}

func NewReimbursementRepository(
	db *postgres.Postgres,
) ReimbursementRepositoryProvider {
	return &reimbursementRepository{
		db: newDBRepo(
			db,
		),
	}
}

func (r *reimbursementRepository) InsertReimbursement(ctx context.Context, rmb reimbursement.Reimbursement) (int, error) {
	id, err := r.db.InsertReimbursement(ctx, rmb)
	if err != nil {
		return 0, err
	}
	return id, nil
}
