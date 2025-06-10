package attendance

const (
	queryInsertAttendacePeriod = `
		INSERT INTO attendance_periods (
			start_date,
			end_date
		) VALUES (
		 	$1,
			$2
		) RETURNING id;
	`
	queryGetAttendancePeriodByID = `
		SELECT 
			id,
			start_date,
			end_date,
			created_at,
			updated_at
		FROM attendance_periods
		WHERE id = $1;
	`

	queryInsertAttendance = `
		INSERT INTO attendances (
			user_id,
			period_id,
			date
		) VALUES (
			$1,
			$2,
			$3
		) RETURNING id;
	`
	queryGetAttendance = `
		SELECT 
			id,
			user_id,
			period_id,
			date,
			created_at,
			updated_at
		FROM attendances
		WHERE user_id = $1 AND period_id = $2 AND date = $3;
	`

	queryGetEmployeeAttendanceSummary = `
		WITH attendance_count AS (
		SELECT user_id, COUNT(*) AS present_days
		FROM attendances
		WHERE period_id = $1
		GROUP BY user_id
		),
		overtime_sum AS (
		SELECT user_id, COALESCE(SUM(hours), 0) AS overtime_hours
		FROM overtimes
		WHERE period_id = $1
		GROUP BY user_id
		),
		reimbursement_sum AS (
		SELECT user_id, COALESCE(SUM(amount), 0) AS reimbursement_total
		FROM reimbursements
		WHERE period_id = $1
		GROUP BY user_id
		)
		SELECT
		u.id AS user_id,
		u.salary AS base_salary,
		COALESCE(a.present_days, 0) AS present_days,
		COALESCE(o.overtime_hours, 0) AS overtime_hours,
		COALESCE(r.reimbursement_total, 0) AS reimbursement_total
		FROM users u
		LEFT JOIN attendance_count a ON a.user_id = u.id
		LEFT JOIN overtime_sum o ON o.user_id = u.id
		LEFT JOIN reimbursement_sum r ON r.user_id = u.id
		WHERE u.is_admin = false
		ORDER BY u.id;
		`
)