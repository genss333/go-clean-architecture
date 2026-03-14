CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    department_id INTEGER NOT NULL,
    hourly_rate_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_employees_departments
        FOREIGN KEY (department_id)
        REFERENCES departments(id) ON DELETE CASCADE,
    CONSTRAINT fk_employees_hourly_rate
        FOREIGN KEY (hourly_rate_id)
        REFERENCES houry_rate(id) ON DELETE CASCADE
);
