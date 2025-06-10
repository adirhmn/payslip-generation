package auth

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	FullName     string `json:"full_name"`
	Salary       int    `json:"salary"`
	IsAdmin      bool   `json:"is_admin"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}