package user

import (
	"context"
	"database/sql"
	"payslip-generation-system/internal/common/errors"
	usermodel "payslip-generation-system/internal/entity/user"
	"payslip-generation-system/internal/postgres"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_newDBRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error preparing mock: %s", err)
	}
	defer db.Close()

	type args struct {
		dbJulo *postgres.Postgres
	}
	tests := []struct {
		name string
		args args
		want dbRepoProvider
	}{
		{
			name: "Happy Path",
			args: args{
				dbJulo: &postgres.Postgres{
					DB: db,
				},
			},
			want: &dbRepo{
				db: &postgres.Postgres{
					DB: db,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDBRepo(tt.args.dbJulo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDBRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dbRepo_GetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockUser := getMockUser()

	type fields struct {
		db *postgres.Postgres
	}
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		mock    func()
		args    args
		want    usermodel.User
		wantErr bool
	}{
		{
			name:   "Happy Path - User Found",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				rows := getMockUserRows(mockUser)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserByUsername)).
					WithArgs(mockUser.Username).
					WillReturnRows(rows)
			},
			args: args{
				ctx:      context.Background(),
				username: mockUser.Username,
			},
			want: usermodel.User{
				ID:           mockUser.ID,
				Username:     mockUser.Username,
				PasswordHash: mockUser.PasswordHash,
				IsAdmin:      mockUser.IsAdmin,
			},
			wantErr: false,
		},
		{
			name:   "User Not Found",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserByUsername)).
					WithArgs("nouser").
					WillReturnError(sql.ErrNoRows)
			},
			args: args{
				ctx:      context.Background(),
				username: "nouser",
			},
			want:    usermodel.User{},
			wantErr: false, 
		},
		{
			name:   "Database Error",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetUserByUsername)).
					WithArgs(mockUser.Username).
					WillReturnError(errors.New("connection error"))
			},
			args: args{
				ctx:      context.Background(),
				username: mockUser.Username,
			},
			want:    usermodel.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			r := &dbRepo{
				db: tt.fields.db,
			}
			got, err := r.GetUserByUsername(tt.args.ctx, tt.args.username)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func getMockUser() usermodel.User {
	return usermodel.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "$2a$10$abcdefghijklmnopqrstuv", 
		IsAdmin:      false,
	}
}

func getMockUserRows(user usermodel.User) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "username", "password_hash", "is_admin"}).
		AddRow(user.ID, user.Username, user.PasswordHash, user.IsAdmin)
}