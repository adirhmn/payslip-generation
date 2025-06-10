package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"payslip-generation-system/internal/entity/audit"
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

func Test_dbRepo_InsertRequestLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockLog := getMockRequestLog()

	type fields struct {
		db *postgres.Postgres
	}
	type args struct {
		ctx context.Context
		log audit.RequestLog
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
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertRequestLog)).
					WithArgs(mockLog.URL, mockLog.Method, mockLog.IPAddress).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockLog.ID))
			},
			args: args{
				ctx: context.Background(),
				log: audit.RequestLog{
					URL:       mockLog.URL,
					Method:    mockLog.Method,
					IPAddress: mockLog.IPAddress,
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name:   "Error Insert",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertRequestLog)).
					WithArgs(mockLog.URL, mockLog.Method, mockLog.IPAddress).
					WillReturnError(sql.ErrConnDone)
			},
			args: args{
				ctx: context.Background(),
				log: audit.RequestLog{
					URL:       mockLog.URL,
					Method:    mockLog.Method,
					IPAddress: mockLog.IPAddress,
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
			got, err := r.InsertRequestLog(tt.args.ctx, tt.args.log)
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

func getMockRequestLog() audit.RequestLog {
	return audit.RequestLog{
		ID:        1,
		URL:       "/api/v1/login",
		Method:    "POST",
		IPAddress: "192.168.1.1",
		CreatedAt: time.Now(),
	}
}

func Test_dbRepo_InsertAuditLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockLog, err := getMockAuditLog()
	assert.NoError(t, err, "Failed to create mock audit log")

	type fields struct {
		db *postgres.Postgres
	}
	type args struct {
		ctx context.Context
		log audit.AuditLog
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
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertAuditLog)).
					WithArgs(
						mockLog.TableName,
						mockLog.RecordID,
						mockLog.Action,
						mockLog.OldData,
						mockLog.NewData,
						mockLog.ChangedBy,
						mockLog.RequestID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockLog.ID))
			},
			args: args{
				ctx: context.Background(),
				log: mockLog,
			},
			want:    1,
			wantErr: false,
		},
		{
			name:   "Error Insert",
			fields: fields{db: &postgres.Postgres{DB: db}},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertAuditLog)).
					WithArgs(
						mockLog.TableName,
						mockLog.RecordID,
						mockLog.Action,
						mockLog.OldData,
						mockLog.NewData,
						mockLog.ChangedBy,
						mockLog.RequestID,
					).
					WillReturnError(sql.ErrConnDone)
			},
			args: args{
				ctx: context.Background(),
				log: mockLog,
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
			got, err := r.InsertAuditLog(tt.args.ctx, tt.args.log)
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

func getMockAuditLog() (audit.AuditLog, error) {
	oldPayload := struct {
		Salary int `json:"salary"`
	}{Salary: 5000000}
	oldData, err := json.Marshal(oldPayload)
	if err != nil {
		return audit.AuditLog{}, err
	}

	newPayload := struct {
		Salary int `json:"salary"`
	}{Salary: 5500000}
	newData, err := json.Marshal(newPayload)
	if err != nil {
		return audit.AuditLog{}, err
	}

	return audit.AuditLog{
		ID:        1,
		TableName: "users",
		RecordID:  101,
		Action:    "UPDATE",
		OldData:   oldData,
		NewData:   newData,
		ChangedBy: sql.NullInt32{Int32: 1, Valid: true},
		RequestID: sql.NullInt32{Int32: 999, Valid: true},
		CreatedAt: time.Now(),
	}, nil
}