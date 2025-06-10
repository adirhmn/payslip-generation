package user

const (
	queryGetUserByUsername = `
		SELECT 
			id, 
			username,
			password_hash,
			is_admin 
		FROM users 
		WHERE username=$1
		ORDER BY created_at DESC
		LIMIT 1
	`

	queryGetAllEmployees = `
		SELECT 
			id,
			username,
			password_hash,
			full_name,
			salary,
			is_admin,
			created_at,
			updated_at
		FROM users
		WHERE is_admin = false;
	`
)