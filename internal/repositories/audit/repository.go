package audit

import (
	"context"

	"payslip-generation-system/internal/entity/audit"
	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository.go -package=mock -destination=mock/repository_mock.go
type AuditRepositoryProvider interface {
	InsertRequestLog(ctx context.Context, log audit.RequestLog) (int, error) 
	InsertAuditLog(ctx context.Context, log audit.AuditLog) (int, error)
}

type auditRepository struct {
	db    dbRepoProvider
}

func NewAuditRepository(
	db *postgres.Postgres,
) AuditRepositoryProvider {
	return &auditRepository{
		db: newDBRepo(
			db,
		),
	}
}

func (r *auditRepository) InsertRequestLog(ctx context.Context, log audit.RequestLog) (int, error) {
	id, err := r.db.InsertRequestLog(ctx, log)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *auditRepository) InsertAuditLog(ctx context.Context, log audit.AuditLog) (int, error) {
	id, err := r.db.InsertAuditLog(ctx, log)
	if err != nil {
		return 0, err
	}
	return id, nil
}
