CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    salary INTEGER NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

CREATE TABLE IF NOT EXISTS attendance_periods (
    id SERIAL PRIMARY KEY,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_attendance_periods_start_end ON attendance_periods(start_date, end_date);

CREATE TABLE IF NOT EXISTS attendances (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    period_id INT NOT NULL REFERENCES attendance_periods(id),
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_attendances_user_id ON attendances(user_id);
CREATE INDEX IF NOT EXISTS idx_attendances_period_id ON attendances(period_id);
CREATE INDEX IF NOT EXISTS idx_attendances_date ON attendances(date);


CREATE TABLE IF NOT EXISTS overtimes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    period_id INT NOT NULL REFERENCES attendance_periods(id),
    date DATE NOT NULL,
    hours INT CHECK (hours >= 1 AND hours <= 3),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_overtimes_user_id ON overtimes(user_id);
CREATE INDEX IF NOT EXISTS idx_overtimes_period_id ON overtimes(period_id);
CREATE INDEX IF NOT EXISTS idx_overtimes_date ON overtimes(date);

CREATE TABLE IF NOT EXISTS reimbursements (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    period_id INT NOT NULL REFERENCES attendance_periods(id),
    amount INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reimbursements_user_id ON reimbursements(user_id);
CREATE INDEX IF NOT EXISTS idx_reimbursements_period_id ON reimbursements(period_id);

CREATE TABLE IF NOT EXISTS payslips (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    period_id INT NOT NULL REFERENCES attendance_periods(id),
    base_salary INTEGER NOT NULL,
    working_days INTEGER NOT NULL,
    present_days INTEGER NOT NULL,
    attendance_amount INTEGER NOT NULL,
    overtime_hours INTEGER NOT NULL,
    overtime_amount INTEGER NOT NULL,
    reimbursement_total INTEGER NOT NULL,
    take_home_pay INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_payslips_user_id ON payslips(user_id);
CREATE INDEX IF NOT EXISTS idx_payslips_period_id ON payslips(period_id);

CREATE TABLE IF NOT EXISTS request_logs (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    method VARCHAR(100) NOT NULL,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(100) NOT NULL,
    record_id INT NOT NULL,
    action VARCHAR(10) NOT NULL CHECK (action IN ('CREATE', 'UPDATE', 'DELETE')),
    old_data JSONB,
    new_data JSONB,
    changed_by INT REFERENCES users(id),
    request_id INT,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_table_record ON audit_logs(table_name, record_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_changed_by ON audit_logs(changed_by);
CREATE INDEX IF NOT EXISTS idx_audit_logs_request_id ON audit_logs(request_id);