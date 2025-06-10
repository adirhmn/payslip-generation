package attendance

import (
	"context"
	"database/sql"
	"payslip-generation-system/internal/entity/attendance"
	"payslip-generation-system/internal/postgres"
	"time"
)

//go:generate mockgen -source=repository_db.go -package=mock -destination=mock/repository_db_mock.go
type dbRepoProvider interface {
	InsertAttendancePeriod(ctx context.Context, attendancePeriod attendance.AttendancePeriod) (int, error)
	GetAttendancePeriodByID(ctx context.Context, id int) (attendance.AttendancePeriod, error) 
	InsertAttendance(ctx context.Context, attendance attendance.Attendance) (int, error) 
	GetAttendance(ctx context.Context, userID, periodID int, date time.Time) (attendance.Attendance, error)
	GetEmployeeAttendanceSummary(ctx context.Context, periodID int) ([]attendance.EmployeeAttendanceSummary, error)
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

func (r *dbRepo) InsertAttendancePeriod(ctx context.Context, attendancePeriod attendance.AttendancePeriod) (int, error){
    var id int
    err := r.db.DB.QueryRowContext(
		ctx,
		queryInsertAttendacePeriod,
		attendancePeriod.StartDate, 
		attendancePeriod.EndDate,
	).Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

func (r *dbRepo) GetAttendancePeriodByID(ctx context.Context, id int) (attendance.AttendancePeriod, error) {
	row := r.db.DB.QueryRowContext(ctx, queryGetAttendancePeriodByID, id)

	var ap attendance.AttendancePeriod

	err := row.Scan(&ap.ID, &ap.StartDate, &ap.EndDate, &ap.CreatedAt, &ap.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return attendance.AttendancePeriod{}, nil
		}
		return attendance.AttendancePeriod{}, err
	}

	return ap, nil
}

func (r *dbRepo) InsertAttendance(ctx context.Context, attendance attendance.Attendance) (int, error) {
	var id int
	err := r.db.DB.QueryRowContext(
		ctx,
		queryInsertAttendance,
		attendance.UserID,
		attendance.PeriodID,
		attendance.Date,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *dbRepo) GetAttendance(ctx context.Context, userID, periodID int, date time.Time) (attendance.Attendance, error) {
	row := r.db.DB.QueryRowContext(ctx, queryGetAttendance, userID, periodID, date)

	var a attendance.Attendance
	err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.PeriodID,
		&a.Date,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return attendance.Attendance{}, nil // or custom error if needed
		}
		return attendance.Attendance{}, err
	}
	return a, nil
}

func (r *dbRepo) GetEmployeeAttendanceSummary(ctx context.Context, periodID int) ([]attendance.EmployeeAttendanceSummary, error) {
    rows, err := r.db.DB.QueryContext(ctx, queryGetEmployeeAttendanceSummary, periodID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []attendance.EmployeeAttendanceSummary
    for rows.Next() {
        var eas attendance.EmployeeAttendanceSummary
        if err := rows.Scan(
            &eas.UserID,
            &eas.BaseSalary,
            &eas.PresentDays,
            &eas.OvertimeHours,
            &eas.ReimbursementTotal,
        ); err != nil {
            return nil, err
        }
        results = append(results, eas)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }
    return results, nil
}

