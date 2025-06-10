package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"payslip-generation-system/internal/entity/attendance"
	"payslip-generation-system/internal/entity/audit"
	"payslip-generation-system/internal/entity/payslip"
	attrepo "payslip-generation-system/internal/repositories/attendance"
	ovttrepo "payslip-generation-system/internal/repositories/overtime"
	payrepo "payslip-generation-system/internal/repositories/payslip"
	rmbrepo "payslip-generation-system/internal/repositories/reimbursement"
	userepo "payslip-generation-system/internal/repositories/user"
	audsvc "payslip-generation-system/internal/services/audit"
)

//go:generate mockgen -source=service.go -package=mock -destination=mock/service_mock.go
type AdminServiceProvider interface {
    AddPeriod(ctx context.Context, attendancePeriod attendance.AttendancePeriod, userID, requestID int)(int, error) 
    RunPayroll(ctx context.Context, periodID, userID, requestID int)( error) 
    GetPayslipSummary(ctx context.Context, periodID int)(payslip.PayslipSummaryReport, error)
}

type adminService struct {
    attrepo attrepo.AttendanceRepositoryProvider
    payrepo payrepo.PayslipRepositoryProvider
    rmbrepo rmbrepo.ReimbursementRepositoryProvider
    ovtrepo ovttrepo.OvertimeRepositoryProvider
    userepo userepo.UserRepositoryProvider
    audsvc audsvc.AuditServiceProvider
}

func NewAdminService(
    attendanceRepo attrepo.AttendanceRepositoryProvider,
    payslipRepo payrepo.PayslipRepositoryProvider,
    reimbursRepo rmbrepo.ReimbursementRepositoryProvider,
    overtimeRepo ovttrepo.OvertimeRepositoryProvider,
    userRepo userepo.UserRepositoryProvider,
    auditService audsvc.AuditServiceProvider,
) AdminServiceProvider {
    return &adminService{
        attrepo: attendanceRepo,
        payrepo: payslipRepo,
        rmbrepo: reimbursRepo,
        ovtrepo: overtimeRepo,
        userepo: userRepo,
        audsvc: auditService,
    }
}

func (s *adminService) AddPeriod(ctx context.Context, attendancePeriod attendance.AttendancePeriod, userID, requestID int)(int, error)  {
    // validaton date
    if !attendancePeriod.StartDate.Before(attendancePeriod.EndDate) {
        return 0, fmt.Errorf("start_date must be before end_date")
    }

    id, err := s.attrepo.InsertAttendancePeriod(ctx, attendancePeriod)
    if err != nil{
        return 0 , err
    }

    attendancePeriodJson, err := json.Marshal(attendancePeriod)
    if err != nil {
        return 0, err
    }

    log := audit.AuditLog{
        TableName: "attendance_periods",
        RecordID: id,
        Action: "CREATE",
        OldData: []byte("{}"),
        NewData: attendancePeriodJson,
        ChangedBy: sql.NullInt32{Valid: true, Int32: int32(userID)},
        RequestID: sql.NullInt32{Valid: true, Int32: int32(requestID)},
    }
    _, err= s.audsvc.RecordAuditLog(ctx, log)
    if err != nil {
        return 0, err
    }
    return id, nil
}

func (s *adminService) RunPayroll(ctx context.Context, periodID, userID, requestID int)( error)  {
    isExistPeriod, err := s.payrepo.PayslipExistsByPeriodID(ctx, periodID)
    if err != nil {
        return  err
    }

    if isExistPeriod {
        return  fmt.Errorf("payroll already generated")
    }

    attendancePeriod,err:=s.attrepo.GetAttendancePeriodByID(ctx, periodID)
    if err != nil {
        return err
    }
    if attendancePeriod.ID == 0 {
        return fmt.Errorf("period not found")
    }
    startDate := attendancePeriod.StartDate
    endDate := attendancePeriod.EndDate
    diff := endDate.Sub(startDate)
    workingDays := int(diff.Hours() / 24) + 1

    employeeSummaries, err := s.attrepo.GetEmployeeAttendanceSummary(ctx, periodID)
    if err != nil {
        return  err
    }

    payslips := []payslip.Payslip{}
    for _, employee := range employeeSummaries {
        attendanceAmount := int((employee.PresentDays*employee.BaseSalary) / workingDays)
        overtimeAmount := int((employee.OvertimeHours * employee.BaseSalary) / workingDays)
        takeHomePay := attendanceAmount + overtimeAmount + employee.ReimbursementTotal

        payslip := payslip.Payslip{
            UserID: employee.UserID,
            PeriodID: periodID,
            BaseSalary: employee.BaseSalary,
            WorkingDays: workingDays,
            PresentDays: employee.PresentDays,
            AttendanceAmount: attendanceAmount,
            OvertimeHours: employee.OvertimeHours,
            OvertimeAmount: overtimeAmount,
            ReimbursementTotal: employee.ReimbursementTotal,
            TakeHomePay: takeHomePay,
        }
        payslips = append(payslips, payslip)
    }

    err = s.payrepo.BulkInsertPayslips(ctx, payslips)
    if err != nil {
        return  err
    }

    payslipsJson, err := json.Marshal(payslips)
    if err != nil {
        return err
    }

    log := audit.AuditLog{
        TableName: "payslips",
        RecordID: 0,
        Action: "CREATE",
        OldData: []byte("{}"),
        NewData: payslipsJson,
        ChangedBy: sql.NullInt32{Valid: true, Int32: int32(userID)},
        RequestID: sql.NullInt32{Valid: true, Int32: int32(requestID)},
    }
    _, err= s.audsvc.RecordAuditLog(ctx, log)
    if err != nil {
        return err
    }

    return nil
}

func (s *adminService) GetPayslipSummary(ctx context.Context, periodID int)(payslip.PayslipSummaryReport, error)  {
   return s.payrepo.GetPayslipSummary(ctx, periodID)
}
