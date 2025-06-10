package payslip

type Payslip struct {
	ID                 int
	UserID             int
	PeriodID           int
	BaseSalary         int
	WorkingDays        int
	PresentDays        int
	AttendanceAmount   int
	OvertimeHours      int
	OvertimeAmount     int
	ReimbursementTotal int
	TakeHomePay        int
	CreatedAt          string
	UpdatedAt          string
}

type PayslipSummary struct {
	UserID        int
	TotalTakeHome int
}

type PayslipSummaryReport struct {
	PerUser []PayslipSummary
	Total   int
}