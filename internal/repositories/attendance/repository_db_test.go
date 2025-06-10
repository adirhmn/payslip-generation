package attendance

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"payslip-generation-system/internal/entity/attendance"
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

func Test_dbRepo_InsertAttendancePeriod(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error preparing mock: %s", err)
	}
	defer db.Close()

	mocktimenow := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	mockAttendancePeriod := getMockAttendancePeriod(mocktimenow)

	type fields struct {
		db *postgres.Postgres
	}
	type args struct {
		ctx  context.Context
		data attendance.AttendancePeriod
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
			name: "Happy Path",
			fields: fields{
				db: &postgres.Postgres{
					DB: db,
				},
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertAttendacePeriod)).
					WithArgs(
						mockAttendancePeriod.StartDate,
						mockAttendancePeriod.EndDate,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(mockAttendancePeriod.ID))
			},
			args: args{
				ctx: context.Background(),
				data: attendance.AttendancePeriod{
					StartDate: mockAttendancePeriod.StartDate,
					EndDate:   mockAttendancePeriod.EndDate,
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Error Insert",
			fields: fields{
				db: &postgres.Postgres{
					DB: db,
				},
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertAttendacePeriod)).
					WithArgs(
						mockAttendancePeriod.StartDate,
						mockAttendancePeriod.EndDate,
					).
					WillReturnError(sql.ErrConnDone)
			},
			args: args{
				ctx: context.Background(),
				data: attendance.AttendancePeriod{
					StartDate: mockAttendancePeriod.StartDate,
					EndDate:   mockAttendancePeriod.EndDate,
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			r := &dbRepo{
				db: tt.fields.db,
			}
			got, err := r.InsertAttendancePeriod(tt.args.ctx, tt.args.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got, "InsertAttendancePeriod = %v, want %v", got, tt.want)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got, "InsertAttendancePeriod = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dbRepo_GetAttendancePeriodByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error preparing mock: %s", err)
	}
	defer db.Close()
	mocktimenow := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	type fields struct {
		db  *postgres.Postgres
	}
	type args struct {
		ctx           context.Context
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		mock    func()
		args    args
		want    attendance.AttendancePeriod
		wantErr bool
	}{
		{
			name: "Happy Path",
			fields: fields{
				db: &postgres.Postgres{
					DB: db,
				},
			},
			mock: func() {
				expectedRows := getMockAttendancePeriodExpectedRows(mocktimenow)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAttendancePeriodByID)).
					WithArgs(1).
					WillReturnRows(expectedRows)
			},
			args: args{
				ctx:           context.Background(),
				id: 1,
			},
			want:    getMockAttendancePeriod(mocktimenow),
			wantErr: false,
		},
		{
			name: "Error - no rows",
			fields: fields{
				db: &postgres.Postgres{
					DB: db,
				},
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAttendancePeriodByID)).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			args: args{
				ctx:           context.Background(),
				id: 1,
			},
			want:    attendance.AttendancePeriod{},
			wantErr: false,
		},
		{
			name: "Error - insert",
			fields: fields{
				db: &postgres.Postgres{
					DB: db,
				},
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAttendancePeriodByID)).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			args: args{
				ctx:           context.Background(),
				id: 1,
			},
			want:    attendance.AttendancePeriod{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt.mock()
		t.Run(tt.name, func(t *testing.T) {
			r := &dbRepo{
				db: tt.fields.db,
			}
			got, err := r.GetAttendancePeriodByID(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got, "GetAttendancePeriodByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getMockAttendancePeriod(mocktime time.Time)  attendance.AttendancePeriod {
	return attendance.AttendancePeriod{
		ID:          1,
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
		IsProcessed: false,
		CreatedAt:   mocktime,
		UpdatedAt:   mocktime,
	}
}

func getMockAttendancePeriodExpectedRows(mocktime time.Time) *sqlmock.Rows {
	mockAttendancePeriod := getMockAttendancePeriod(mocktime)
	expectedRowsData := [][]driver.Value{
		{
			mockAttendancePeriod.ID,
			mockAttendancePeriod.StartDate,
			mockAttendancePeriod.EndDate,
			mockAttendancePeriod.CreatedAt,
			mockAttendancePeriod.UpdatedAt,
		},
	}

	expectedRows := sqlmock.NewRows([]string{
		"id",
		"start_date",
		"end_date",
		"created_at",
		"updated_at",
	})

	for _, row := range expectedRowsData {
		expectedRows.AddRow(row...)
	}

	return expectedRows
}


func Test_dbRepo_InsertAttendance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error preparing mock: %s", err)
	}
	defer db.Close()

	mocktimenow := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	mockAttendance := getMockAttendance(mocktimenow)

	type fields struct {
		db *postgres.Postgres
	}
	type args struct {
		ctx  context.Context
		data attendance.Attendance
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
			name: "Happy Path",
			fields: fields{
				db: &postgres.Postgres{DB: db},
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertAttendance)).
					WithArgs(
						mockAttendance.UserID,
						mockAttendance.PeriodID,
						mockAttendance.Date,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockAttendance.ID))
			},
			args: args{
				ctx: context.Background(),
				data: attendance.Attendance{
					UserID:   mockAttendance.UserID,
					PeriodID: mockAttendance.PeriodID,
					Date:     mockAttendance.Date,
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Error Insert",
			fields: fields{
				db: &postgres.Postgres{DB: db},
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertAttendance)).
					WithArgs(
						mockAttendance.UserID,
						mockAttendance.PeriodID,
						mockAttendance.Date,
					).
					WillReturnError(sql.ErrConnDone)
			},
			args: args{
				ctx: context.Background(),
				data: attendance.Attendance{
					UserID:   mockAttendance.UserID,
					PeriodID: mockAttendance.PeriodID,
					Date:     mockAttendance.Date,
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
			got, err := r.InsertAttendance(tt.args.ctx, tt.args.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got, "InsertAttendance() got = %v, want %v", got, tt.want)
		})
	}
}

func Test_dbRepo_GetAttendance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error preparing mock: %s", err)
	}
	defer db.Close()
	mocktimenow := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	mockAttendance := getMockAttendance(mocktimenow)

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
		want    attendance.Attendance
		wantErr bool
	}{
		{
			name: "Happy Path",
			fields: fields{
				db: &postgres.Postgres{DB: db},
			},
			mock: func() {
				expectedRows := getMockAttendanceExpectedRows(mocktimenow)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAttendance)).
					WithArgs(mockAttendance.UserID, mockAttendance.PeriodID, mockAttendance.Date).
					WillReturnRows(expectedRows)
			},
			args: args{
				ctx:      context.Background(),
				userID:   mockAttendance.UserID,
				periodID: mockAttendance.PeriodID,
				date:     mockAttendance.Date,
			},
			want:    mockAttendance,
			wantErr: false,
		},
		{
			name: "Error - no rows",
			fields: fields{
				db: &postgres.Postgres{DB: db},
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAttendance)).
					WithArgs(mockAttendance.UserID, mockAttendance.PeriodID, mockAttendance.Date).
					WillReturnError(sql.ErrNoRows)
			},
			args: args{
				ctx:      context.Background(),
				userID:   mockAttendance.UserID,
				periodID: mockAttendance.PeriodID,
				date:     mockAttendance.Date,
			},
			want:    attendance.Attendance{},
			wantErr: false,
		},
		{
			name: "Error - database",
			fields: fields{
				db: &postgres.Postgres{DB: db},
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetAttendance)).
					WithArgs(mockAttendance.UserID, mockAttendance.PeriodID, mockAttendance.Date).
					WillReturnError(sql.ErrConnDone)
			},
			args: args{
				ctx:      context.Background(),
				userID:   mockAttendance.UserID,
				periodID: mockAttendance.PeriodID,
				date:     mockAttendance.Date,
			},
			want:    attendance.Attendance{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			r := &dbRepo{
				db: tt.fields.db,
			}
			got, err := r.GetAttendance(tt.args.ctx, tt.args.userID, tt.args.periodID, tt.args.date)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got, "GetAttendance() got = %v, want %v", got, tt.want)
		})
	}
}

func getMockAttendance(mocktime time.Time) attendance.Attendance {
	return attendance.Attendance{
		ID:        1,
		UserID:    101,
		PeriodID:  202,
		Date:      time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC),
		CreatedAt: mocktime,
		UpdatedAt: mocktime,
	}
}


func getMockAttendanceExpectedRows(mocktime time.Time) *sqlmock.Rows {
	mockAtt := getMockAttendance(mocktime)
	
	rows := sqlmock.NewRows([]string{
		"id",
		"user_id",
		"period_id",
		"date",
		"created_at",
		"updated_at",
	})

	rows.AddRow(
		mockAtt.ID,
		mockAtt.UserID,
		mockAtt.PeriodID,
		mockAtt.Date,
		mockAtt.CreatedAt,
		mockAtt.UpdatedAt,
	)

	return rows
}

func Test_dbRepo_GetEmployeeAttendanceSummary(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockPeriodID := 123
	mockData := getMockEmployeeAttendanceSummaryData()

	type fields struct {
		db *postgres.Postgres
	}
	type args struct {
		ctx      context.Context
		periodID int
	}
	tests := []struct {
		name    string
		fields  fields
		mock    func()
		args    args
		want    []attendance.EmployeeAttendanceSummary
		wantErr bool
	}{
		{
			name:   "Happy Path - Multiple Rows",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				rows := getMockEmployeeAttendanceSummaryRows(mockData)
				mock.ExpectQuery(regexp.QuoteMeta(queryGetEmployeeAttendanceSummary)).
					WithArgs(mockPeriodID).
					WillReturnRows(rows)
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID},
			want:    mockData,
			wantErr: false,
		},
		{
			name:   "Happy Path - No Rows",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "base_salary", "present_days", "overtime_hours", "reimbursement_total"})
				mock.ExpectQuery(regexp.QuoteMeta(queryGetEmployeeAttendanceSummary)).
					WithArgs(mockPeriodID).
					WillReturnRows(rows)
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID},
			want:    []attendance.EmployeeAttendanceSummary{}, // Expecting an empty slice, not nil
			wantErr: false,
		},
		{
			name:   "Error - Query Failed",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryGetEmployeeAttendanceSummary)).
					WithArgs(mockPeriodID).
					WillReturnError(sql.ErrConnDone)
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "Error - Scan Failed",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "base_salary", "present_days", "overtime_hours", "reimbursement_total"}).
					AddRow(mockData[0].UserID, mockData[0].BaseSalary, mockData[0].PresentDays, mockData[0].OvertimeHours, mockData[0].ReimbursementTotal).
					AddRow("invalid_user_id", "invalid_salary", "invalid_days", "invalid_hours", "invalid_reimbursement") // Bad data to cause scan error

				mock.ExpectQuery(regexp.QuoteMeta(queryGetEmployeeAttendanceSummary)).
					WithArgs(mockPeriodID).
					WillReturnRows(rows)
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			r := &dbRepo{
				db: tt.fields.db,
			}
			got, err := r.GetEmployeeAttendanceSummary(tt.args.ctx, tt.args.periodID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.ElementsMatch(t, tt.want, got)
			
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

// getMockEmployeeAttendanceSummaryData creates mock data for the summary.
func getMockEmployeeAttendanceSummaryData() []attendance.EmployeeAttendanceSummary {
	return []attendance.EmployeeAttendanceSummary{
		{
			UserID:             101,
			BaseSalary:         5000000,
			PresentDays:        20,
			OvertimeHours:      3, 
			ReimbursementTotal: 150000,
		},
		{
			UserID:             102,
			BaseSalary:         7500000,
			PresentDays:        22,
			OvertimeHours:      3, 
			ReimbursementTotal: 50000,
		},
	}
}

func getMockEmployeeAttendanceSummaryRows(data []attendance.EmployeeAttendanceSummary) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"user_id",
		"base_salary",
		"present_days",
		"overtime_hours",
		"reimbursement_total",
	})
	for _, item := range data {
		rows.AddRow(item.UserID, item.BaseSalary, item.PresentDays, item.OvertimeHours, item.ReimbursementTotal)
	}
	return rows
}