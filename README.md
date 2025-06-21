# Payslip Generation System API

The system ensures fairness and automation in payroll processing and enables employees to transparently track their monthly earnings. Main features:
1. Allow employees to submit attendance, overtime, and reimbursement requests securely and reliably.
2. Enable admins to generate accurate payslips for all employees based on locked attendance and claims.


---

## Features

- Employee attendance submission (excluding weekends)
- Overtime submission (up to 3 hours/day)
- Reimbursement requests with descriptions
- Admin payroll period management and payroll generation
- Payslip generation and summary reports for employees and admin
- Role-based authentication (Admin & Employee)
- One-time payroll run per payroll period (freezes data)

---

## API Endpoints

### Authentication

All endpoints require **Basic Auth** headers. Admin routes require admin privileges.

---

### Employee APIs (`/private/employee`)

| Endpoint                      | Method | Description                            |
|-------------------------------|--------|------------------------------------|
| `/attendance/submit`          | POST   | Submit daily attendance (no weekends) |
| `/overtime/submit`            | POST   | Submit overtime hours (max 3/day)  |
| `/reimbursement/submit`       | POST   | Submit reimbursement request        |
| `/payslips/:period_id`        | GET    | Get payslip breakdown for a payroll period |

---

### Admin APIs (`/private/admin`)

| Endpoint                                 | Method | Description                           |
|------------------------------------------|--------|-----------------------------------|
| `/generate-payroll/:period_id`           | POST   | Run payroll process for given period (locks data) |
| `/payslips/:period_id`                   | GET    | Get summary of all employee payslips for a period |
| `/payslips/:period_id/:user_id`          | GET    | Get payslip breakdown for specific employee |

---

## Setup & Run
1. Clone repository
2. Setup your database and configure connection
3. Run migrations / seed initial data (users + payroll_periods)
4. Start the server

