package overtime

import (
	"context"
	"database/sql"
	"payslip-generation-system/internal/entity/overtime"
	"payslip-generation-system/internal/postgres"
	"time"
)

//go:generate mockgen -source=repository_db.go -package=mock -destination=mock/repository_db_mock.go
type dbRepoProvider interface {
	InsertOvertime(ctx context.Context, ot overtime.Overtime) (int, error) 
	GetOvertime(ctx context.Context, userID, periodID int, date time.Time) (overtime.Overtime, error)
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

func (r *dbRepo) InsertOvertime(ctx context.Context, ot overtime.Overtime) (int, error) {
	var id int
	err := r.db.DB.QueryRowContext(
		ctx,
		queryInsertOvertime,
		ot.UserID,
		ot.PeriodID,
		ot.Date,
		ot.Hours,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *dbRepo) GetOvertime(ctx context.Context, userID, periodID int, date time.Time) (overtime.Overtime, error) {
	row := r.db.DB.QueryRowContext(ctx, queryGetOvertime, userID, periodID, date)

	var ot overtime.Overtime
	err := row.Scan(
		&ot.ID,
		&ot.UserID,
		&ot.PeriodID,
		&ot.Date,
		&ot.Hours,
		&ot.CreatedAt,
		&ot.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return overtime.Overtime{}, nil
		}
		return overtime.Overtime{}, err
	}

	return ot, nil
}
