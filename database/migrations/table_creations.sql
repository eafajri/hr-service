-- ENUM types
CREATE TYPE user_role AS ENUM ('employee', 'admin');
CREATE TYPE payroll_periods_status AS ENUM ('open', 'closed');

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    salary NUMERIC(10, 2) DEFAULT 0,
    role user_role NOT NULL
);

-- Attendance table
CREATE TABLE attendances (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    check_in_time TIMESTAMP,
    check_out_time TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, date)
);

-- Overtime table
CREATE TABLE overtimes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    durations INT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, date)
);

-- Reimbursements table
CREATE TABLE reimbursements (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    amount NUMERIC(10, 2) DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- Payroll periods table
CREATE TABLE payroll_periods (
    id SERIAL PRIMARY KEY,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    status payroll_periods_status DEFAULT 'open',
    UNIQUE(period_start, period_end)
);

-- Payslips table
CREATE TABLE payslips (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    payroll_period_id INT NOT NULL REFERENCES payroll_periods(id) ON DELETE CASCADE,
    base_salary NUMERIC(10, 2) NOT NULL,
    attendance_days INT NOT NULL,
    overtime_hours INT NOT NULL,
    overtime_pay NUMERIC(10, 2) NOT NULL,
    reimbursement_total NUMERIC(10, 2) NOT NULL,
    total_take_home NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, payroll_period_id)
);
