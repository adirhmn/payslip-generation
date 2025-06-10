package overtime

import "time"

type Overtime struct {
	ID        int
	UserID    int
	PeriodID  int
	Date      time.Time 
	Hours     int
	CreatedAt time.Time
	UpdatedAt time.Time
}