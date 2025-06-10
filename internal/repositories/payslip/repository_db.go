package payslip

import (
	"context"
	"fmt"
	"payslip-generation-system/internal/entity/payslip"
	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository_db.go -package=mock -destination=mock/repository_db_mock.go
type dbRepoProvider interface {
	BulkInsertPayslips(ctx context.Context, payslips []payslip.Payslip) error 
	PayslipExistsByPeriodID(ctx context.Context, periodID int) (bool, error)
	GetPayslipsByUserID(ctx context.Context, userID int) ([]payslip.Payslip, error)
	GetPayslipSummary(ctx context.Context, periodID int) (payslip.PayslipSummaryReport, error)
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

func (r *dbRepo) BulkInsertPayslips(ctx context.Context, payslips []payslip.Payslip) error {
	if len(payslips) == 0 {
		return nil
	}

	query := queryBulkInsertPayslips

	args := []interface{}{}
	values := ""

	for i, p := range payslips {
		start := i * 10
		values += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			start+1, start+2, start+3, start+4, start+5,
			start+6, start+7, start+8, start+9, start+10,
		)
		args = append(args,
			p.UserID, p.PeriodID, p.BaseSalary, p.WorkingDays, p.PresentDays,
			p.AttendanceAmount, p.OvertimeHours, p.OvertimeAmount, p.ReimbursementTotal, p.TakeHomePay,
		)
	}

	// Remove trailing comma
	query = query + values[:len(values)-1] + ";"

	_, err := r.db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *dbRepo) PayslipExistsByPeriodID(ctx context.Context, periodID int) (bool, error) {
	var exists bool
	err := r.db.DB.QueryRowContext(ctx, queryCheckPayslipExistsByPeriodID, periodID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *dbRepo) GetPayslipsByUserID(ctx context.Context, userID int) ([]payslip.Payslip, error) {
	rows, err := r.db.DB.QueryContext(ctx, queryGetPayslipsByUserID, userID)
	if err != nil {
		return []payslip.Payslip{}, err
	}
	defer rows.Close()

	var payslips []payslip.Payslip
	for rows.Next() {
		var p payslip.Payslip
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.PeriodID,
			&p.BaseSalary,
			&p.WorkingDays,
			&p.PresentDays,
			&p.AttendanceAmount,
			&p.OvertimeHours,
			&p.OvertimeAmount,
			&p.ReimbursementTotal,
			&p.TakeHomePay,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return []payslip.Payslip{}, err
		}
		payslips = append(payslips, p)
	}

	if err = rows.Err(); err != nil {
		return []payslip.Payslip{}, err
	}

	return payslips, nil
}

func (r *dbRepo) GetPayslipSummary(ctx context.Context, periodID int) (payslip.PayslipSummaryReport, error) {
	rows, err := r.db.DB.QueryContext(ctx, queryPayslipSummaryPerUser, periodID)
	if err != nil {
		return payslip.PayslipSummaryReport{}, err
	}
	defer rows.Close()

	var summaries []payslip.PayslipSummary
	for rows.Next() {
		var s payslip.PayslipSummary
		if err := rows.Scan(&s.UserID, &s.TotalTakeHome); err != nil {
			return payslip.PayslipSummaryReport{}, err
		}
		summaries = append(summaries, s)
	}
	if err := rows.Err(); err != nil {
		return payslip.PayslipSummaryReport{}, err
	}

	var total int
	err = r.db.DB.QueryRowContext(ctx, queryPayslipSummaryTotal, periodID).Scan(&total)
	if err != nil {
		return payslip.PayslipSummaryReport{}, err
	}

	return payslip.PayslipSummaryReport{
		PerUser: summaries,
		Total:   total,
	}, nil
}