package attendance

import "time"

type AttendancePeriod struct {
	ID          int32
	StartDate   time.Time
	EndDate     time.Time
	IsProcessed bool
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
