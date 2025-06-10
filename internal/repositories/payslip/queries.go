package payslip

const (
	queryCheckPayslipExistsByPeriodID = `
		SELECT EXISTS (
			SELECT 1 FROM payslips WHERE period_id = $1
		);
	`

	queryBulkInsertPayslips = `
		INSERT INTO payslips (
				user_id, period_id, base_salary, working_days, present_days,
				attendance_amount, overtime_hours, overtime_amount, reimbursement_total, take_home_pay
			) VALUES 
		`

	queryGetPayslipsByUserID = `
		SELECT
			id,
			user_id,
			period_id,
			base_salary,
			working_days,
			present_days,
			attendance_amount,
			overtime_hours,
			overtime_amount,
			reimbursement_total,
			take_home_pay,
			created_at,
			updated_at
		FROM payslips
		WHERE user_id = $1;
		`

	queryPayslipSummaryPerUser = `
		SELECT user_id, SUM(take_home_pay) AS total_take_home
		FROM payslips
		WHERE period_id = $1
		GROUP BY user_id;
	`

	queryPayslipSummaryTotal = `
		SELECT COALESCE(SUM(take_home_pay), 0) AS total_take_home
		FROM payslips
		WHERE period_id = $1;
	`
)
