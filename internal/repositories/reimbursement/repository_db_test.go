package reimbursement

import (
	"context"
	"database/sql"
	"payslip-generation-system/internal/entity/reimbursement"
	"payslip-generation-system/internal/postgres"
	"reflect"
	"regexp"
	"testing"
	"time"

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

func Test_dbRepo_InsertReimbursement(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error preparing mock: %s", err)
	}
	defer db.Close()

	mockTimeNow := time.Now()
	mockRmb := getMockReimbursement(mockTimeNow)

	type fields struct {
		db *postgres.Postgres
	}
	type args struct {
		ctx context.Context
		rmb reimbursement.Reimbursement
	}
	tests := []struct {
		name    string
		fields  fields
		mock    func()
		args    args
		want    int
		wantErr bool
	}{
		{
			name:   "Happy Path",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertReimbursement)).
					WithArgs(mockRmb.UserID, mockRmb.PeriodID, mockRmb.Amount, mockRmb.Description).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockRmb.ID))
			},
			args: args{
				ctx: context.Background(),
				rmb: reimbursement.Reimbursement{
					UserID:      mockRmb.UserID,
					PeriodID:    mockRmb.PeriodID,
					Amount:      mockRmb.Amount,
					Description: mockRmb.Description,
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name:   "Error Insert",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertReimbursement)).
					WithArgs(mockRmb.UserID, mockRmb.PeriodID, mockRmb.Amount, mockRmb.Description).
					WillReturnError(sql.ErrConnDone)
			},
			args: args{
				ctx: context.Background(),
				rmb: reimbursement.Reimbursement{
					UserID:      mockRmb.UserID,
					PeriodID:    mockRmb.PeriodID,
					Amount:      mockRmb.Amount,
					Description: mockRmb.Description,
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			r := &dbRepo{
				db: tt.fields.db,
			}
			got, err := r.InsertReimbursement(tt.args.ctx, tt.args.rmb)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)

			// Memastikan semua ekspektasi mock terpenuhi
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}


func getMockReimbursement(mocktime time.Time) reimbursement.Reimbursement {
	return reimbursement.Reimbursement{
		ID:          1,
		UserID:      101,
		PeriodID:    202406,
		Amount:      150000,
		Description: "Biaya transport",
		CreatedAt:   mocktime,
		UpdatedAt:   mocktime,
	}
}