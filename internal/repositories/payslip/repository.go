package payslip

import (
	"context"

	"payslip-generation-system/internal/entity/payslip"
	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository.go -package=mock -destination=mock/repository_mock.go
type PayslipRepositoryProvider interface {
	BulkInsertPayslips(ctx context.Context, payslips []payslip.Payslip) error
	PayslipExistsByPeriodID(ctx context.Context, periodID int) (bool, error)
	GetPayslipsByUserID(ctx context.Context, userID int) ([]payslip.Payslip, error) 
	GetPayslipSummary(ctx context.Context, periodID int) (payslip.PayslipSummaryReport, error)
}

type payslipRepository struct {
	db    dbRepoProvider
}

func NewPayslipRepository(
	db *postgres.Postgres,
) PayslipRepositoryProvider {
	return &payslipRepository{
		db: newDBRepo(
			db,
		),
	}
}

func (r *payslipRepository) BulkInsertPayslips(ctx context.Context, payslips []payslip.Payslip) error {
	err := r.db.BulkInsertPayslips(ctx, payslips)
	if err != nil {
		return err
	}
	return nil
}

func (r *payslipRepository) PayslipExistsByPeriodID(ctx context.Context, periodID int) (bool, error) {
	exists, err := r.db.PayslipExistsByPeriodID(ctx, periodID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *payslipRepository) GetPayslipsByUserID(ctx context.Context, userID int) ([]payslip.Payslip, error) {
	payslips, err := r.db.GetPayslipsByUserID(ctx, userID)
	if err != nil {
		return []payslip.Payslip{}, err
	}
	return payslips, nil
}

func (r *payslipRepository) GetPayslipSummary(ctx context.Context, periodID int) (payslip.PayslipSummaryReport, error) {
	report, err := r.db.GetPayslipSummary(ctx, periodID)
	if err != nil {
		return payslip.PayslipSummaryReport{}, err
	}
	return report, nil
}