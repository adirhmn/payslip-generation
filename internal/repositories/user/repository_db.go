package user

import (
	"context"
	"database/sql"
	"payslip-generation-system/internal/common/errors"
	usermodel "payslip-generation-system/internal/entity/user"
	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository_db.go -package=mock -destination=mock/repository_db_mock.go
type dbRepoProvider interface {
	GetUserByUsername(ctx context.Context, username string) (usermodel.User, error)
	GetAllEmployees(ctx context.Context) ([]usermodel.User, error) 
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

func (r *dbRepo)  GetUserByUsername(ctx context.Context, username string) (usermodel.User, error) {
    var u usermodel.User
    err := r.db.DB.QueryRowContext(ctx, queryGetUserByUsername, username).Scan(
		&u.ID, 
		&u.Username,
		&u.PasswordHash,
		&u.IsAdmin,
	)
    if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usermodel.User{}, nil
		}
		return usermodel.User{}, err
	}
    return u, nil
}

func (r *dbRepo) GetAllEmployees(ctx context.Context) ([]usermodel.User, error) {
	rows, err := r.db.DB.QueryContext(ctx, queryGetAllEmployees)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []usermodel.User
	for rows.Next() {
		var u usermodel.User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.PasswordHash,
			&u.FullName,
			&u.Salary,
			&u.IsAdmin,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}
