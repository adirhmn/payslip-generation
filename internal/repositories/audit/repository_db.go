package audit

import (
	"context"
	"payslip-generation-system/internal/entity/audit"
	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository_db.go -package=mock -destination=mock/repository_db_mock.go
type dbRepoProvider interface {
	InsertRequestLog(ctx context.Context, log audit.RequestLog) (int, error)
	InsertAuditLog(ctx context.Context, log audit.AuditLog) (int, error) 
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

func (r *dbRepo) InsertRequestLog(ctx context.Context, log audit.RequestLog) (int, error) {
	var id int
	err := r.db.DB.QueryRowContext(
		ctx,
		queryInsertRequestLog,
		log.URL,
		log.Method,
		log.IPAddress,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *dbRepo) InsertAuditLog(ctx context.Context, log audit.AuditLog) (int, error) {
	var id int
	err := r.db.DB.QueryRowContext(
		ctx,
		queryInsertAuditLog,
		log.TableName,
		log.RecordID,
		log.Action,
		log.OldData,
		log.NewData,
		log.ChangedBy,
		log.RequestID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
