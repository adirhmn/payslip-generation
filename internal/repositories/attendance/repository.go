package attendance

import (
	"context"
	"time"

	"payslip-generation-system/internal/entity/attendance"
	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository.go -package=mock -destination=mock/repository_mock.go
type AttendanceRepositoryProvider interface {
	InsertAttendancePeriod(ctx context.Context, attendancePeriod attendance.AttendancePeriod) (int, error)
	GetAttendancePeriodByID(ctx context.Context, id int) (attendance.AttendancePeriod, error)
	InsertAttendance(ctx context.Context, a attendance.Attendance) (int, error)
	GetAttendance(ctx context.Context, userID, periodID int, date time.Time) (attendance.Attendance, error)
	GetEmployeeAttendanceSummary(ctx context.Context, periodID int) ([]attendance.EmployeeAttendanceSummary, error)
}

type attendanceRepository struct {
	db    dbRepoProvider
}

func NewAttendanceRepository(
	db *postgres.Postgres,
) AttendanceRepositoryProvider {
	return &attendanceRepository{
		db: newDBRepo(
			db,
		),
	}
}


func (r *attendanceRepository) InsertAttendancePeriod(ctx context.Context, attendancePeriod attendance.AttendancePeriod) (int, error){
	id, err := r.db.InsertAttendancePeriod(ctx, attendancePeriod)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *attendanceRepository) GetAttendancePeriodByID(ctx context.Context, id int) (attendance.AttendancePeriod, error) {
	result, err := r.db.GetAttendancePeriodByID(ctx, id)
	if err != nil {
		return attendance.AttendancePeriod{}, err
	}
	return result, nil
}

func (r *attendanceRepository) InsertAttendance(ctx context.Context, a attendance.Attendance) (int, error) {
	id, err := r.db.InsertAttendance(ctx, a)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *attendanceRepository) GetAttendance(ctx context.Context, userID, periodID int, date time.Time) (attendance.Attendance, error) {
	result, err := r.db.GetAttendance(ctx, userID, periodID, date)
	if err != nil {
		return attendance.Attendance{}, err
	}

	return result, nil
}

func (r *attendanceRepository) GetEmployeeAttendanceSummary(ctx context.Context, periodID int) ([]attendance.EmployeeAttendanceSummary, error) {
    result, err := r.db.GetEmployeeAttendanceSummary(ctx, periodID)
    if err != nil {
        return []attendance.EmployeeAttendanceSummary{}, err
    }

    return result, nil
}
