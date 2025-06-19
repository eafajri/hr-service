# Context

In a company, there is data that contains employees' salaries. They're paid with the same rule, which is monthly-based, with regular 8 working hours per day (9AM-5PM), 5 days a week (monday-friday). Their take-home pay will be prorated based on their attendance. Along with that, they can also propose overtime, which is paid at twice the prorated salary for hours taken. They can also submit reimbursement requests which will be included in the payslip.

# Features

As employee:
- submit attendance
- submit overtime
- submit reimbursements
- generate payslip

As admin
- run payroll (close the period)
- generate employee payroll details

# First setup
1. Run docker
2. Install posgres
```
docker pull postgres
```
3. Create a new database
```
docker run --name postgres_db \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=hris_db \
  -p 5432:5432 \
  -d postgres
```
4. Execute migration on cmd/migrations