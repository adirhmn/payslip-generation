package user

import (
	"context"

	usermodel "payslip-generation-system/internal/entity/user"
	"payslip-generation-system/internal/postgres"
)

//go:generate mockgen -source=repository.go -package=mock -destination=mock/repository_mock.go
type UserRepositoryProvider interface {
	GetUserByUsername(ctx context.Context,username string) (usermodel.User, error)
	GetAllEmployees(ctx context.Context) ([]usermodel.User, error)
}

type userRepository struct {
	db    dbRepoProvider
}

func NewUserRepository(
	db *postgres.Postgres,
) UserRepositoryProvider {
	return &userRepository{
		db: newDBRepo(
			db,
		),
	}
}

func (r *userRepository) GetUserByUsername(ctx context.Context,username string) (usermodel.User, error){
	// ping (all) database
	user, err := r.db.GetUserByUsername(ctx, username)
	if err != nil {
		return usermodel.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetAllEmployees(ctx context.Context) ([]usermodel.User, error) {
	employees, err := r.db.GetAllEmployees(ctx)
	if err != nil {
		return nil, err
	}
	return employees, nil
}