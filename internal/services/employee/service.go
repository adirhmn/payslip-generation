package employee

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"payslip-generation-system/internal/entity/attendance"
	"payslip-generation-system/internal/entity/audit"
	"payslip-generation-system/internal/entity/overtime"
	"payslip-generation-system/internal/entity/payslip"
	"payslip-generation-system/internal/entity/reimbursement"
	attrepo "payslip-generation-system/internal/repositories/attendance"
	ovtrepo "payslip-generation-system/internal/repositories/overtime"
	payreporepo "payslip-generation-system/internal/repositories/payslip"
	rmbrepo "payslip-generation-system/internal/repositories/reimbursement"
	audsvc "payslip-generation-system/internal/services/audit"
)

//go:generate mockgen -source=service.go -package=mock -destination=mock/service_mock.go
type EmployeeServiceProvider interface {
    SubmitAttendance(ctx context.Context, attendance attendance.Attendance, requestID int)(int, error) 
	SubmitOvertime(ctx context.Context, overtime overtime.Overtime, requestID int)(int, error) 
	SubmitReimbursement(ctx context.Context, reimbursement reimbursement.Reimbursement, requestID int)(int, error)
	GeneratePayslips(ctx context.Context, userID int)([]payslip.Payslip, error)
}

type employeeService struct {
    attrepo attrepo.AttendanceRepositoryProvider
	ovtrepo ovtrepo.OvertimeRepositoryProvider
	rmbrepo rmbrepo.ReimbursementRepositoryProvider
	payrepo payreporepo.PayslipRepositoryProvider
    audsvc audsvc.AuditServiceProvider
}

func NewEmployeeService(
    attendanceRepo attrepo.AttendanceRepositoryProvider,
	overtimeRepo ovtrepo.OvertimeRepositoryProvider,
	reimbursRepo rmbrepo.ReimbursementRepositoryProvider,
	payslipRepo payreporepo.PayslipRepositoryProvider,
    auditService audsvc.AuditServiceProvider,
) EmployeeServiceProvider {
    return &employeeService{
        attrepo: attendanceRepo,
		ovtrepo: overtimeRepo,
		rmbrepo: reimbursRepo,
		payrepo: payslipRepo,
        audsvc: auditService,
    }
}



func (s *employeeService) SubmitAttendance(ctx context.Context, attendance attendance.Attendance, requestID int)(int, error) {
    existingAttendance , err:= s.attrepo.GetAttendance(ctx, attendance.UserID, attendance.PeriodID, attendance.Date)
    if err != nil {
        return 0, err
    }

    if existingAttendance.ID != 0 {
        return 0, fmt.Errorf("attendance already exists")
    }
    
    attendancePeriod, err:= s.attrepo.GetAttendancePeriodByID(ctx, attendance.PeriodID)
    if err != nil {
        return 0, err
    }
    if attendancePeriod.ID == 0{
        return 0, fmt.Errorf("period not found")
    }

    if attendance.Date.Before(attendancePeriod.StartDate) || attendance.Date.After(attendancePeriod.EndDate) {
        return 0, fmt.Errorf("date must be between %s and %s", attendancePeriod.StartDate, attendancePeriod.EndDate)
    }

    id, err := s.attrepo.InsertAttendance(ctx, attendance)
    if err != nil {
        return 0, err
    }

    attendanceJson, err := json.Marshal(attendance)
    if err != nil {
        return 0, err
    }

    log := audit.AuditLog{
        TableName: "attendances",
        RecordID: id,
        Action: "CREATE",
        OldData: []byte("{}"),
        NewData: attendanceJson,
        ChangedBy: sql.NullInt32{Valid: true, Int32: int32(attendance.UserID)},
        RequestID: sql.NullInt32{Valid: true, Int32: int32(requestID)},
    }
    _, err= s.audsvc.RecordAuditLog(ctx, log)
    if err != nil {
        return 0, err
    }
    return id, nil
}

func (s *employeeService) SubmitOvertime(ctx context.Context, overtime overtime.Overtime, requestID int)(int, error) {
	existingAttendance , err:= s.attrepo.GetAttendance(ctx, overtime.UserID, overtime.PeriodID, overtime.Date)
    if err != nil {
        return 0, err
    }

    if existingAttendance.ID == 0 {
        return 0, fmt.Errorf("you need to submit attendance first before submitting overtime")
    }

	existingOvertime , err:= s.ovtrepo.GetOvertime(ctx, overtime.UserID, overtime.PeriodID, overtime.Date)
    if err != nil {
        return 0, err
    }

	if existingOvertime.ID != 0 {
		return 0, fmt.Errorf("overtime already exists")
	}

	attendancePeriod, err:= s.attrepo.GetAttendancePeriodByID(ctx, overtime.PeriodID)
    if err != nil {
        return 0, err
    }
    if attendancePeriod.ID == 0{
        return 0, fmt.Errorf("period not found")
    }

	if overtime.Hours > 3 || overtime.Hours < 1 {
		return 0, fmt.Errorf("hours must be between 1 and 3")
	}

    s.ovtrepo.InsertOvertime(ctx, overtime)
    id, err := s.ovtrepo.InsertOvertime(ctx, overtime)
    if err != nil {
        return 0, err
    }

    overtimeJson, err := json.Marshal(overtime)
    if err != nil {
        return 0, err
    }

    log := audit.AuditLog{
        TableName: "overtimes",
        RecordID: id,
        Action: "CREATE",
        OldData: []byte("{}"),
        NewData: overtimeJson,
        ChangedBy: sql.NullInt32{Valid: true, Int32: int32(overtime.UserID)},
        RequestID: sql.NullInt32{Valid: true, Int32: int32(requestID)},
    }
    _, err= s.audsvc.RecordAuditLog(ctx, log)
    if err != nil {
        return 0, err
    }

	return id, nil
}

func (s *employeeService) SubmitReimbursement(ctx context.Context, reimbursement reimbursement.Reimbursement, requestID int)(int, error) {
	attendancePeriod, err:= s.attrepo.GetAttendancePeriodByID(ctx, reimbursement.PeriodID)
    if err != nil {
        return 0, err
    }
    if attendancePeriod.ID == 0{
        return 0, fmt.Errorf("period not found")
    }

    id, err := s.rmbrepo.InsertReimbursement(ctx, reimbursement)
    if err != nil {
        return 0, err
    }

    reimbursementJson, err := json.Marshal(reimbursement)
    if err != nil {
        return 0, err
    }

    log := audit.AuditLog{
        TableName: "reimbursements",
        RecordID: id,
        Action: "CREATE",
        OldData: []byte("{}"),
        NewData: reimbursementJson,
        ChangedBy: sql.NullInt32{Valid: true, Int32: int32(reimbursement.UserID)},
        RequestID: sql.NullInt32{Valid: true, Int32: int32(requestID)},
    }
    _, err= s.audsvc.RecordAuditLog(ctx, log)
    if err != nil {
        return 0, err
    }
	return id, nil
}

func (s *employeeService) GeneratePayslips(ctx context.Context, userID int)([]payslip.Payslip, error) {
	return s.payrepo.GetPayslipsByUserID(ctx, userID)
}