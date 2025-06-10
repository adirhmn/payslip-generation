package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"payslip-generation-system/internal/entity/attendance"
	"payslip-generation-system/internal/entity/audit"
	"payslip-generation-system/internal/entity/payslip"
	attrepo "payslip-generation-system/internal/repositories/attendance"
	mockattrepo "payslip-generation-system/internal/repositories/attendance/mock"
	ovttrepo "payslip-generation-system/internal/repositories/overtime"
	mockovttrepo "payslip-generation-system/internal/repositories/overtime/mock"
	payrepo "payslip-generation-system/internal/repositories/payslip"
	mockpayrepo "payslip-generation-system/internal/repositories/payslip/mock"
	rmbrepo "payslip-generation-system/internal/repositories/reimbursement"
	mockrmbrepo "payslip-generation-system/internal/repositories/reimbursement/mock"
	userepo "payslip-generation-system/internal/repositories/user"
	mockuserepo "payslip-generation-system/internal/repositories/user/mock"
	audsvc "payslip-generation-system/internal/services/audit"
	mockaudsvc "payslip-generation-system/internal/services/audit/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewAdminService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPayRepo := mockpayrepo.NewMockPayslipRepositoryProvider(ctrl)
	mockOvtRepo:=mockovttrepo.NewMockOvertimeRepositoryProvider(ctrl)
	mockAttRepo:=mockattrepo.NewMockAttendanceRepositoryProvider(ctrl)
	mockRmbRepo :=mockrmbrepo.NewMockReimbursementRepositoryProvider(ctrl)
	mockAudSvc :=mockaudsvc.NewMockAuditServiceProvider(ctrl)
	mockUserRepo:= mockuserepo.NewMockUserRepositoryProvider(ctrl)

	type args struct {
		attrepo attrepo.AttendanceRepositoryProvider
		payrepo payrepo.PayslipRepositoryProvider
		rmbrepo rmbrepo.ReimbursementRepositoryProvider
		ovtrepo ovttrepo.OvertimeRepositoryProvider
		userepo userepo.UserRepositoryProvider
		audsvc audsvc.AuditServiceProvider
	}
	tests := []struct {
		name string
		args args
		want AdminServiceProvider
	}{
		{
			name: "Happy Path",
			args: args{
				attrepo: mockAttRepo,
				payrepo: mockPayRepo,
				rmbrepo: mockRmbRepo,
				ovtrepo: mockOvtRepo,
				userepo: mockUserRepo,
				audsvc: mockAudSvc,
			},
			want: &adminService{
				attrepo: mockAttRepo,
				payrepo: mockPayRepo,
				rmbrepo: mockRmbRepo,
				ovtrepo: mockOvtRepo,
				userepo: mockUserRepo,
				audsvc: mockAudSvc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAdminService(tt.args.attrepo, tt.args.payrepo, tt.args.rmbrepo, tt.args.ovtrepo, tt.args.userepo, tt.args.audsvc)
			assert.Equal(t, got, tt.want, "NewAdminService() = %v, want %v", got, tt.want)
		})
	}
}

func Test_adminService_AddPeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAttRepo := mockattrepo.NewMockAttendanceRepositoryProvider(ctrl)
	mockAudSvc := mockaudsvc.NewMockAuditServiceProvider(ctrl)

	mockUserID := 1
	mockRequestID := 99
	mockStartDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	mockEndDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	validPeriod := attendance.AttendancePeriod{
		StartDate: mockStartDate,
		EndDate:   mockEndDate,
	}
	
	validPeriodJSON, _ := json.Marshal(validPeriod)

	type args struct {
		ctx              context.Context
		attendancePeriod attendance.AttendancePeriod
		userID           int
		requestID        int
	}
	tests := []struct {
		name    string
		mock    func()
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Happy Path - Success",
			mock: func() {
				mockAttRepo.EXPECT().
					InsertAttendancePeriod(gomock.Any(), validPeriod).
					Return(1, nil).
					Times(1)

				expectedLog := audit.AuditLog{
					TableName: "attendance_periods",
					RecordID:  1, 
					Action:    "CREATE",
					OldData:   []byte("{}"),
					NewData:   validPeriodJSON,
					ChangedBy: sql.NullInt32{Valid: true, Int32: int32(mockUserID)},
					RequestID: sql.NullInt32{Valid: true, Int32: int32(mockRequestID)},
				}
				mockAudSvc.EXPECT().
					RecordAuditLog(gomock.Any(), gomock.Eq(expectedLog)).
					Return(1, nil).
					Times(1)
			},
			args: args{
				ctx:              context.Background(),
				attendancePeriod: validPeriod,
				userID:           mockUserID,
				requestID:        mockRequestID,
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Error - Invalid Date",
			mock: func() {
				
			},
			args: args{
				ctx: context.Background(),
				attendancePeriod: attendance.AttendancePeriod{
					StartDate: mockEndDate,   
					EndDate:   mockStartDate,
				},
				userID:    mockUserID,
				requestID: mockRequestID,
			},
			want:    0,
			wantErr: true, 
		},
		{
			name: "Error - InsertAttendancePeriod failed",
			mock: func() {
				mockAttRepo.EXPECT().
					InsertAttendancePeriod(gomock.Any(), validPeriod).
					Return(0, errors.New("database connection error")).
					Times(1)
			},
			args: args{
				ctx:              context.Background(),
				attendancePeriod: validPeriod,
				userID:           mockUserID,
				requestID:        mockRequestID,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Error - RecordAuditLog failed",
			mock: func() {
				mockAttRepo.EXPECT().
					InsertAttendancePeriod(gomock.Any(), validPeriod).
					Return(1, nil).
					Times(1)
				
				mockAudSvc.EXPECT().
					RecordAuditLog(gomock.Any(), gomock.Any()). 
					Return(0, errors.New("audit service down")).
					Times(1)
			},
			args: args{
				ctx:              context.Background(),
				attendancePeriod: validPeriod,
				userID:           mockUserID,
				requestID:        mockRequestID,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			s := NewAdminService(mockAttRepo, nil, nil, nil, nil, mockAudSvc)

			got, err := s.AddPeriod(tt.args.ctx, tt.args.attendancePeriod, tt.args.userID, tt.args.requestID)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got, "AddPeriod() = %v, want %v", got, tt.want)
		})
	}
}

func Test_adminService_RunPayroll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPayRepo := mockpayrepo.NewMockPayslipRepositoryProvider(ctrl)
	mockAttRepo := mockattrepo.NewMockAttendanceRepositoryProvider(ctrl)
	mockAudSvc := mockaudsvc.NewMockAuditServiceProvider(ctrl)

	mockPeriodID := 202506
	mockUserID := 1
	mockRequestID := 101

	mockStartDate := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	mockEndDate := time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)
	mockWorkingDays := 10
	mockPeriod := attendance.AttendancePeriod{
		ID:        int32(mockPeriodID),
		StartDate: mockStartDate,
		EndDate:   mockEndDate,
	}

	mockSummaries := []attendance.EmployeeAttendanceSummary{
		{UserID: 10, BaseSalary: 3000000, PresentDays: 8, OvertimeHours: 60, ReimbursementTotal: 100000},
	}

	// expected calculation
	emp := mockSummaries[0]
	expectedAttendanceAmount := (emp.PresentDays * emp.BaseSalary) / mockWorkingDays      
	expectedOvertimeAmount := (emp.OvertimeHours * emp.BaseSalary) / mockWorkingDays   
	expectedTakeHomePay := expectedAttendanceAmount + expectedOvertimeAmount + emp.ReimbursementTotal 

	expectedPayslips := []payslip.Payslip{
		{
			UserID:             emp.UserID,
			PeriodID:           mockPeriodID,
			BaseSalary:         emp.BaseSalary,
			WorkingDays:        mockWorkingDays,
			PresentDays:        emp.PresentDays,
			AttendanceAmount:   expectedAttendanceAmount,
			OvertimeHours:      emp.OvertimeHours,
			OvertimeAmount:     expectedOvertimeAmount,
			ReimbursementTotal: emp.ReimbursementTotal,
			TakeHomePay:        expectedTakeHomePay,
		},
	}
	expectedPayslipsJSON, _ := json.Marshal(expectedPayslips)
	expectedAuditLog := audit.AuditLog{
		TableName: "payslips", RecordID: 0, Action: "CREATE", OldData: []byte("{}"), NewData: expectedPayslipsJSON,
		ChangedBy: sql.NullInt32{Valid: true, Int32: int32(mockUserID)},
		RequestID: sql.NullInt32{Valid: true, Int32: int32(mockRequestID)},
	}

	type args struct {
		ctx       context.Context
		periodID  int
		userID    int
		requestID int
	}
	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy Path",
			mock: func() {
				gomock.InOrder(
					mockPayRepo.EXPECT().PayslipExistsByPeriodID(gomock.Any(), mockPeriodID).Return(false, nil),
					mockAttRepo.EXPECT().GetAttendancePeriodByID(gomock.Any(), mockPeriodID).Return(mockPeriod, nil),
					mockAttRepo.EXPECT().GetEmployeeAttendanceSummary(gomock.Any(), mockPeriodID).Return(mockSummaries, nil),
					mockPayRepo.EXPECT().BulkInsertPayslips(gomock.Any(), gomock.Eq(expectedPayslips)).Return(nil),
					mockAudSvc.EXPECT().RecordAuditLog(gomock.Any(), gomock.Eq(expectedAuditLog)).Return(1, nil),
				)
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID, userID: mockUserID, requestID: mockRequestID},
			wantErr: assert.NoError,
		},
		{
			name: "Error - Payroll Already Exists",
			mock: func() {
				mockPayRepo.EXPECT().PayslipExistsByPeriodID(gomock.Any(), mockPeriodID).Return(true, nil)
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID, userID: mockUserID, requestID: mockRequestID},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "payroll already generated")
			},
		},
		{
			name: "Error - GetAttendancePeriodByID failed",
			mock: func() {
				mockPayRepo.EXPECT().PayslipExistsByPeriodID(gomock.Any(), mockPeriodID).Return(false, nil)
				mockAttRepo.EXPECT().GetAttendancePeriodByID(gomock.Any(), mockPeriodID).Return(attendance.AttendancePeriod{}, errors.New("db error"))
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID, userID: mockUserID, requestID: mockRequestID},
			wantErr: assert.Error,
		},
		{
			name: "Error - Period Not Found",
			mock: func() {
				mockPayRepo.EXPECT().PayslipExistsByPeriodID(gomock.Any(), mockPeriodID).Return(false, nil)
				mockAttRepo.EXPECT().GetAttendancePeriodByID(gomock.Any(), mockPeriodID).Return(attendance.AttendancePeriod{ID: 0}, nil)
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID, userID: mockUserID, requestID: mockRequestID},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, err, "period not found")
			},
		},
		{
			name: "Error - GetEmployeeAttendanceSummary failed",
			mock: func() {
				mockPayRepo.EXPECT().PayslipExistsByPeriodID(gomock.Any(), mockPeriodID).Return(false, nil)
				mockAttRepo.EXPECT().GetAttendancePeriodByID(gomock.Any(), mockPeriodID).Return(mockPeriod, nil)
				mockAttRepo.EXPECT().GetEmployeeAttendanceSummary(gomock.Any(), mockPeriodID).Return(nil, errors.New("db error"))
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID, userID: mockUserID, requestID: mockRequestID},
			wantErr: assert.Error,
		},
		{
			name: "Error - BulkInsertPayslips failed",
			mock: func() {
				mockPayRepo.EXPECT().PayslipExistsByPeriodID(gomock.Any(), mockPeriodID).Return(false, nil)
				mockAttRepo.EXPECT().GetAttendancePeriodByID(gomock.Any(), mockPeriodID).Return(mockPeriod, nil)
				mockAttRepo.EXPECT().GetEmployeeAttendanceSummary(gomock.Any(), mockPeriodID).Return(mockSummaries, nil)
				mockPayRepo.EXPECT().BulkInsertPayslips(gomock.Any(), gomock.Eq(expectedPayslips)).Return(errors.New("bulk insert error"))
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID, userID: mockUserID, requestID: mockRequestID},
			wantErr: assert.Error,
		},
		{
			name: "Error - RecordAuditLog failed",
			mock: func() {
				mockPayRepo.EXPECT().PayslipExistsByPeriodID(gomock.Any(), mockPeriodID).Return(false, nil)
				mockAttRepo.EXPECT().GetAttendancePeriodByID(gomock.Any(), mockPeriodID).Return(mockPeriod, nil)
				mockAttRepo.EXPECT().GetEmployeeAttendanceSummary(gomock.Any(), mockPeriodID).Return(mockSummaries, nil)
				mockPayRepo.EXPECT().BulkInsertPayslips(gomock.Any(), gomock.Eq(expectedPayslips)).Return(nil)
				mockAudSvc.EXPECT().RecordAuditLog(gomock.Any(), gomock.Eq(expectedAuditLog)).Return(0, errors.New("audit error"))
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID, userID: mockUserID, requestID: mockRequestID},
			wantErr: assert.Error,
		},
		{
			name: "Edge Case - No Employees",
			mock: func() {
				gomock.InOrder(
					mockPayRepo.EXPECT().PayslipExistsByPeriodID(gomock.Any(), mockPeriodID).Return(false, nil),
					mockAttRepo.EXPECT().GetAttendancePeriodByID(gomock.Any(), mockPeriodID).Return(mockPeriod, nil),

					mockAttRepo.EXPECT().GetEmployeeAttendanceSummary(gomock.Any(), mockPeriodID).Return([]attendance.EmployeeAttendanceSummary{}, nil),

					mockPayRepo.EXPECT().BulkInsertPayslips(gomock.Any(), gomock.Eq([]payslip.Payslip{})).Return(nil),
					mockAudSvc.EXPECT().RecordAuditLog(gomock.Any(), gomock.Any()).Return(1, nil),
				)
			},
			args:    args{ctx: context.Background(), periodID: mockPeriodID, userID: mockUserID, requestID: mockRequestID},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			s := NewAdminService(mockAttRepo, mockPayRepo, nil, nil, nil, mockAudSvc)
			err := s.RunPayroll(tt.args.ctx, tt.args.periodID, tt.args.userID, tt.args.requestID)
			tt.wantErr(t, err)
		})
	}
}

func Test_adminService_GetPayslipSummary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPayRepo := mockpayrepo.NewMockPayslipRepositoryProvider(ctrl)
	mockPeriodID := 202506

	expectedReport := payslip.PayslipSummaryReport{
		PerUser: []payslip.PayslipSummary{
			{UserID: 101, TotalTakeHome: 8500000},
			{UserID: 102, TotalTakeHome: 9200000},
		},
		Total: 2,
	}

	type args struct {
		ctx      context.Context
		periodID int
	}
	tests := []struct {
		name    string
		mock    func()
		args    args
		want    payslip.PayslipSummaryReport
		wantErr bool
	}{
		{
			name: "Happy Path - Success with new struct",
			mock: func() {
				mockPayRepo.EXPECT().
					GetPayslipSummary(gomock.Any(), mockPeriodID).
					Return(expectedReport, nil).
					Times(1)
			},
			args: args{
				ctx:      context.Background(),
				periodID: mockPeriodID,
			},
			want:    expectedReport,
			wantErr: false,
		},
		{
			name: "Error - Repository",
			mock: func() {
				mockPayRepo.EXPECT().
					GetPayslipSummary(gomock.Any(), mockPeriodID).
					Return(payslip.PayslipSummaryReport{}, errors.New("database connection failed")).
					Times(1)
			},
			args: args{
				ctx:      context.Background(),
				periodID: mockPeriodID,
			},
			want:    payslip.PayslipSummaryReport{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			s := NewAdminService(nil, mockPayRepo, nil, nil, nil, nil)

			got, err := s.GetPayslipSummary(tt.args.ctx, tt.args.periodID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}