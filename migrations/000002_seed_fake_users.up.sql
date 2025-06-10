-- Insert 1 admin
INSERT INTO users (username, password_hash, full_name, salary, is_admin, created_at, updated_at)
VALUES 
('admin', '$2a$10$pfUBo.b5KR0y6E.FDrfbX.YkHV0jFqVKeQGheFPK0LUYc5POx9i0y', 'Admin User', 0, TRUE, NOW(), NOW());

-- Insert 100 employees
INSERT INTO users (username, password_hash, full_name, salary, is_admin, created_at, updated_at)
SELECT 
    'employee_' || i,
    '$2a$10$6xNhgdUmzN.MK9pT/gcsNOVZMnv5CdQk/ZaABw8mbYn2PzKPXj9gW',
    'Employee ' || i,
    FLOOR(RANDOM() * 15 + 6) * 500000, -- salary range: 3jt–10,5jt
    FALSE,
    NOW(),
    NOW()
FROM generate_series(1, 100) AS s(i);

