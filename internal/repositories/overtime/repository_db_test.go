package overtime

import (
	"context"
	"database/sql"
	"payslip-generation-system/internal/entity/overtime"
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

func Test_dbRepo_InsertOvertime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error preparing mock: %s", err)
	}
	defer db.Close()

	mocktimenow := time.Now()
	mockOvertime := getMockOvertime(mocktimenow)

	type fields struct {
		db *postgres.Postgres
	}
	type args struct {
		ctx context.Context
		ot  overtime.Overtime
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
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertOvertime)).
					WithArgs(mockOvertime.UserID, mockOvertime.PeriodID, mockOvertime.Date, mockOvertime.Hours).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockOvertime.ID))
			},
			args: args{
				ctx: context.Background(),
				ot: overtime.Overtime{
					UserID:   mockOvertime.UserID,
					PeriodID: mockOvertime.PeriodID,
					Date:     mockOvertime.Date,
					Hours:    mockOvertime.Hours,
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name:   "Error Insert",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertOvertime)).
					WithArgs(mockOvertime.UserID, mockOvertime.PeriodID, mockOvertime.Date, mockOvertime.Hours).
					WillReturnError(sql.ErrConnDone)
			},
			args: args{
				ctx: context.Background(),
				ot: overtime.Overtime{
					UserID:   mockOvertime.UserID,
					PeriodID: mockOvertime.PeriodID,
					Date:     mockOvertime.Date,
					Hours:    mockOvertime.Hours,
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
			got, err := r.InsertOvertime(tt.args.ctx, tt.args.ot)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_dbRepo_GetOvertime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error preparing mock: %s", err)
	}
	defer db.Close()

	mocktimenow := time.Now()
	mockOvertime := getMockOvertime(mocktimenow)

	type fields struct {
		db *postgres.Postgres
	}
	type args struct {
		ctx      context.Context
		userID   int
		periodID int
		date     time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		mock    func()
		args    args
		want    overtime.Overtime
		wantErr bool
	}{
		{
			name:   "Happy Path",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				rows := getMockOvertimeExpectedRows(mocktimenow)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetOvertime)).
					WithArgs(mockOvertime.UserID, mockOvertime.PeriodID, mockOvertime.Date).
					WillReturnRows(rows)
			},
			args: args{
				ctx:      context.Background(),
				userID:   mockOvertime.UserID,
				periodID: mockOvertime.PeriodID,
				date:     mockOvertime.Date,
			},
			want:    mockOvertime,
			wantErr: false,
		},
		{
			name:   "Error - no rows",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetOvertime)).
					WithArgs(mockOvertime.UserID, mockOvertime.PeriodID, mockOvertime.Date).
					WillReturnError(sql.ErrNoRows)
			},
			args: args{
				ctx:      context.Background(),
				userID:   mockOvertime.UserID,
				periodID: mockOvertime.PeriodID,
				date:     mockOvertime.Date,
			},
			want:    overtime.Overtime{},
			wantErr: false,
		},
		{
			name:   "Error - database",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetOvertime)).
					WithArgs(mockOvertime.UserID, mockOvertime.PeriodID, mockOvertime.Date).
					WillReturnError(sql.ErrConnDone)
			},
			args: args{
				ctx:      context.Background(),
				userID:   mockOvertime.UserID,
				periodID: mockOvertime.PeriodID,
				date:     mockOvertime.Date,
			},
			want:    overtime.Overtime{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			r := &dbRepo{
				db: tt.fields.db,
			}
			got, err := r.GetOvertime(tt.args.ctx, tt.args.userID, tt.args.periodID, tt.args.date)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}


func getMockOvertime(mocktime time.Time) overtime.Overtime {
	return overtime.Overtime{
		ID:        1,
		UserID:    101,
		PeriodID:  202,
		Date:      time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
		Hours:     2,
		CreatedAt: mocktime,
		UpdatedAt: mocktime,
	}
}


func getMockOvertimeExpectedRows(mocktime time.Time) *sqlmock.Rows {
	mockOt := getMockOvertime(mocktime)

	rows := sqlmock.NewRows([]string{
		"id",
		"user_id",
		"period_id",
		"date",
		"hours",
		"created_at",
		"updated_at",
	})

	rows.AddRow(
		mockOt.ID,
		mockOt.UserID,
		mockOt.PeriodID,
		mockOt.Date,
		mockOt.Hours,
		mockOt.CreatedAt,
		mockOt.UpdatedAt,
	)

	return rows
}