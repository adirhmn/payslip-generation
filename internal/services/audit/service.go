package audit

import (
	"context"

	"payslip-generation-system/internal/entity/audit"
	repo "payslip-generation-system/internal/repositories/audit"
)

//go:generate mockgen -source=service.go -package=mock -destination=mock/service_mock.go
type AuditServiceProvider interface {
	RecordRequestLog(ctx context.Context, log audit.RequestLog)(int, error) 
	RecordAuditLog(ctx context.Context, log audit.AuditLog)(int, error)
}

type auditService struct {
	repo repo.AuditRepositoryProvider
}

func NewAuditService(
	auditRepo repo.AuditRepositoryProvider,
) AuditServiceProvider {
	return &auditService{
		repo: auditRepo,
	}
}

func (s *auditService) RecordRequestLog(ctx context.Context, log audit.RequestLog)(int, error) {
	return s.repo.InsertRequestLog(ctx, log)
}

func (s *auditService) RecordAuditLog(ctx context.Context, log audit.AuditLog)(int, error) {
	return s.repo.InsertAuditLog(ctx, log)
}