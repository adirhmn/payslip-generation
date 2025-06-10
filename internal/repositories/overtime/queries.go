package overtime

const (
	queryInsertOvertime = `
		INSERT INTO overtimes (
			user_id,
			period_id,
			date,
			hours
		) VALUES (
			$1,
			$2,
			$3,
			$4
		) RETURNING id;
	`

	queryGetOvertime = `
		SELECT 
			id,
			user_id,
			period_id,
			date,
			hours,
			created_at,
			updated_at
		FROM overtimes
		WHERE user_id = $1 AND period_id = $2 AND date = $3;
	`
)