package reimbursement

const (
	queryInsertReimbursement = `
		INSERT INTO reimbursements (
			user_id,
			period_id,
			amount,
			description
		) VALUES (
			$1,
			$2,
			$3,
			$4
		) RETURNING id;
	`
)
