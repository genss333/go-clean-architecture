CREATE TABLE IF NOT EXISTS pay_stubs (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL,
    gross_pay DECIMAL,
    tax DECIMAL,
    net_pay DECIMAL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pay_stubs_employees
        FOREIGN KEY (employee_id)
        REFERENCES employees(id) ON DELETE CASCADE
);
