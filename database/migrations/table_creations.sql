-- ENUM types
CREATE TYPE user_role AS ENUM ('employee', 'admin');
CREATE TYPE payroll_periods_status AS ENUM ('open', 'closed');

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    role user_role NOT NULL
);

CREATE TABLE user_salaries (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(10, 2) NOT NULL,
    effective_from DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Attendance table
CREATE TABLE employee_attendances (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    check_in_time TIMESTAMP NOT NULL,
    check_out_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) NOT NULL,
    UNIQUE (user_id, date)
);

-- Overtime table
CREATE TABLE employee_overtimes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    durations INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) NOT NULL,
    UNIQUE(user_id, date)
);

-- Reimbursements table
CREATE TABLE employee_reimbursements (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    amount NUMERIC(10, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) NOT NULL,
    description TEXT
    UNIQUE (user_id, date)
);

-- Payroll periods table
CREATE TABLE payroll_periods (
    id SERIAL PRIMARY KEY,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    working_days INT NOT NULL,
    status payroll_periods_status DEFAULT 'open',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) NOT NULL,
    UNIQUE(period_start, period_end)
);

-- Payslips table
CREATE TABLE payroll_payslips (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    payroll_period_id INT NOT NULL REFERENCES payroll_periods(id) ON DELETE CASCADE,
    base_salary NUMERIC(10, 2) NOT NULL,
    attendance_days INT NOT NULL,
    attendance_hours INT NOT NULL,
    attendance_pay NUMERIC(10, 2) NOT NULL,
    overtime_hours INT NOT NULL,
    overtime_pay NUMERIC(10, 2) NOT NULL,
    reimbursement_total NUMERIC(10, 2) NOT NULL,
    total_take_home NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    UNIQUE(user_id, payroll_period_id)
);
