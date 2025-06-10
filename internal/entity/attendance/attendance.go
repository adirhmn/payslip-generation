package attendance

import "time"

type Attendance struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PeriodID  int       `json:"period_id"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
