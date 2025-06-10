package reimbursement

import "time"

type Reimbursement struct {
	ID          int
	UserID      int
	PeriodID    int
	Amount      int
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
