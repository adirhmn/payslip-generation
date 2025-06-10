package payslip

import (
	"context"
	"database/sql"
	"payslip-generation-system/internal/entity/payslip"
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

func Test_dbRepo_GetPayslipsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockUserID := 101
	mockData := getMockPayslipsData(mockUserID)

	type fields struct {
		db *postgres.Postgres
	}
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name    string
		fields  fields
		mock    func()
		args    args
		want    []payslip.Payslip
		wantErr bool
	}{
		{
			name:   "Happy Path - Multiple Rows",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				rows := getMockPayslipsRows(mockData)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetPayslipsByUserID)).
					WithArgs(mockUserID).
					WillReturnRows(rows)
			},
			args:    args{ctx: context.Background(), userID: mockUserID},
			want:    mockData,
			wantErr: false,
		},
		{
			name:   "Happy Path - No Rows",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "period_id", "base_salary", "working_days", "present_days", "attendance_amount", "overtime_hours", "overtime_amount", "reimbursement_total", "take_home_pay", "created_at", "updated_at"})
				mock.ExpectQuery(regexp.QuoteMeta(queryGetPayslipsByUserID)).
					WithArgs(mockUserID).
					WillReturnRows(rows)
			},
			args:    args{ctx: context.Background(), userID: mockUserID},
			want:    []payslip.Payslip{}, 
			wantErr: false,
		},
		{
			name:   "Error - Query Failed",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetPayslipsByUserID)).
					WithArgs(mockUserID).
					WillReturnError(sql.ErrConnDone)
			},
			args:    args{ctx: context.Background(), userID: mockUserID},
			want:    []payslip.Payslip{},
			wantErr: true,
		},
		{
			name:   "Error - Scan Failed",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "period_id", "base_salary"}).
					AddRow(1, mockUserID, 202405, "salary is unvalid") 
				mock.ExpectQuery(regexp.QuoteMeta(queryGetPayslipsByUserID)).
					WithArgs(mockUserID).
					WillReturnRows(rows)
			},
			args:    args{ctx: context.Background(), userID: mockUserID},
			want:    []payslip.Payslip{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			r := &dbRepo{
				db: tt.fields.db,
			}
			got, err := r.GetPayslipsByUserID(tt.args.ctx, tt.args.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}


func getMockPayslipsData(userID int) []payslip.Payslip {
	mockTime := time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	return []payslip.Payslip{
		{
			ID:                 1,
			UserID:             userID,
			PeriodID:           202405,
			BaseSalary:         8000000,
			WorkingDays:        22,
			PresentDays:        22,
			AttendanceAmount:   8000000,
			OvertimeHours:      10,
			OvertimeAmount:     500000,
			ReimbursementTotal: 250000,
			TakeHomePay:        8750000,
			CreatedAt:          mockTime,
			UpdatedAt:          mockTime,
		},
		{
			ID:                 2,
			UserID:             userID,
			PeriodID:           202406,
			BaseSalary:         8000000,
			WorkingDays:        20,
			PresentDays:        19,
			AttendanceAmount:   7600000,
			OvertimeHours:      5,
			OvertimeAmount:     250000,
			ReimbursementTotal: 100000,
			TakeHomePay:        7950000,
			CreatedAt:          mockTime,
			UpdatedAt:          mockTime,
		},
	}
}


func getMockPayslipsRows(data []payslip.Payslip) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "period_id", "base_salary", "working_days",
		"present_days", "attendance_amount", "overtime_hours",
		"overtime_amount", "reimbursement_total", "take_home_pay",
		"created_at", "updated_at",
	})
	for _, p := range data {
		rows.AddRow(
			p.ID, p.UserID, p.PeriodID, p.BaseSalary, p.WorkingDays,
			p.PresentDays, p.AttendanceAmount, p.OvertimeHours,
			p.OvertimeAmount, p.ReimbursementTotal, p.TakeHomePay,
			p.CreatedAt, p.UpdatedAt,
		)
	}
	return rows
}