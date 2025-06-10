package audit

import (
	"database/sql"
	"time"
)

type RequestLog struct {
	ID        int
	URL       string
	Method    string
	IPAddress string
	CreatedAt time.Time
}

type AuditLog struct {
	ID         int
	TableName  string
	RecordID   int
	Action     string
	OldData    []byte 
	NewData    []byte
	ChangedBy  sql.NullInt32
	RequestID  sql.NullInt32
	CreatedAt  time.Time
}

