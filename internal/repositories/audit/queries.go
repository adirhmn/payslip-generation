package audit

const (
	queryInsertRequestLog = `
		INSERT INTO request_logs (
			url,
			method,
			ip_address
		) VALUES (
			$1,
			$2,
			$3
		) RETURNING id;
	`

	queryInsertAuditLog = `
		INSERT INTO audit_logs (
			table_name,
			record_id,
			action,
			old_data,
			new_data,
			changed_by,
			request_id
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		) RETURNING id;
	`
)