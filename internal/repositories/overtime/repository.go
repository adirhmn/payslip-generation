package overtime

import (
	"context"
	"time"

	"payslip-generation-system/internal/entity/overtime"
	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository.go -package=mock -destination=mock/repository_mock.go
type OvertimeRepositoryProvider interface {
	InsertOvertime(ctx context.Context, ot overtime.Overtime) (int, error) 
	GetOvertime(ctx context.Context, userID, periodID int, date time.Time) (overtime.Overtime, error)
}

type overtimeRepository struct {
	db    dbRepoProvider
}

func NewOvertimeRepository(
	db *postgres.Postgres,
) OvertimeRepositoryProvider {
	return &overtimeRepository{
		db: newDBRepo(
			db,
		),
	}
}


func (r *overtimeRepository) InsertOvertime(ctx context.Context, ot overtime.Overtime) (int, error) {
	id, err := r.db.InsertOvertime(ctx, ot)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *overtimeRepository) GetOvertime(ctx context.Context, userID, periodID int, date time.Time) (overtime.Overtime, error) {
	result, err := r.db.GetOvertime(ctx, userID, periodID, date)
	if err != nil {
		return overtime.Overtime{}, err
	}
	return result, nil
}
