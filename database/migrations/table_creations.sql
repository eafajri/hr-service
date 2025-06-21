CREATE TYPE user_role AS ENUM ('employee', 'admin');
CREATE TYPE payroll_periods_status AS ENUM ('open', 'closed');

-- public.audit_logs definition

-- Drop table

-- DROP TABLE public.audit_logs;

CREATE TABLE public.audit_logs (
	id serial4 NOT NULL,
	request_id varchar(255) NULL,
	ip_address varchar(64) NULL,
	table_name varchar(255) NOT NULL,
	"action" varchar(255) NOT NULL,
	"target" varchar(255) NOT NULL,
	payload jsonb NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(255) NOT NULL,
	CONSTRAINT audit_logs_pkey PRIMARY KEY (id)
);


-- public.payroll_periods definition

-- Drop table

-- DROP TABLE public.payroll_periods;

CREATE TABLE public.payroll_periods (
	id serial4 NOT NULL,
	period_start date NOT NULL,
	period_end date NOT NULL,
	working_days int4 NOT NULL,
	status public."payroll_periods_status" DEFAULT 'open'::payroll_periods_status NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(255) NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(255) NOT NULL,
	CONSTRAINT payroll_periods_period_start_period_end_key UNIQUE (period_start, period_end),
	CONSTRAINT payroll_periods_pkey PRIMARY KEY (id)
);


-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id serial4 NOT NULL,
	username varchar(255) NOT NULL,
	"password" text NOT NULL,
	"role" public."user_role" NOT NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);


-- public.employee_attendances definition

-- Drop table

-- DROP TABLE public.employee_attendances;

CREATE TABLE public.employee_attendances (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	"date" date NOT NULL,
	check_in_time timestamp NOT NULL,
	check_out_time timestamp NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(255) NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(255) NOT NULL,
	CONSTRAINT employee_attendances_pkey PRIMARY KEY (id),
	CONSTRAINT employee_attendances_user_id_date_key UNIQUE (user_id, date),
	CONSTRAINT employee_attendances_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);


-- public.employee_overtimes definition

-- Drop table

-- DROP TABLE public.employee_overtimes;

CREATE TABLE public.employee_overtimes (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	"date" date NOT NULL,
	durations int4 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(255) NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(255) NOT NULL,
	CONSTRAINT employee_overtimes_pkey PRIMARY KEY (id),
	CONSTRAINT employee_overtimes_user_id_date_key UNIQUE (user_id, date),
	CONSTRAINT employee_overtimes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);


-- public.employee_reimbursements definition

-- Drop table

-- DROP TABLE public.employee_reimbursements;

CREATE TABLE public.employee_reimbursements (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	"date" date NOT NULL,
	amount numeric(10, 2) DEFAULT 0 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(255) NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_by varchar(255) NOT NULL,
	description text NULL,
	CONSTRAINT employee_reimbursements_pkey PRIMARY KEY (id),
	CONSTRAINT unique_user_date UNIQUE (user_id, date),
	CONSTRAINT employee_reimbursements_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);


-- public.payroll_payslips definition

-- Drop table

-- DROP TABLE public.payroll_payslips;

CREATE TABLE public.payroll_payslips (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	payroll_period_id int4 NOT NULL,
	base_salary numeric(10, 2) NOT NULL,
	attendance_days int4 NOT NULL,
	attendance_hours int4 NOT NULL,
	attendance_pay numeric(10, 2) NOT NULL,
	overtime_hours int4 NOT NULL,
	overtime_pay numeric(10, 2) NOT NULL,
	reimbursement_total numeric(10, 2) NOT NULL,
	total_take_home numeric(10, 2) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(255) NOT NULL,
	CONSTRAINT payroll_payslips_pkey PRIMARY KEY (id),
	CONSTRAINT payroll_payslips_user_id_payroll_period_id_key UNIQUE (user_id, payroll_period_id),
	CONSTRAINT payroll_payslips_payroll_period_id_fkey FOREIGN KEY (payroll_period_id) REFERENCES public.payroll_periods(id) ON DELETE CASCADE,
	CONSTRAINT payroll_payslips_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);


-- public.user_salaries definition

-- Drop table

-- DROP TABLE public.user_salaries;

CREATE TABLE public.user_salaries (
	id serial4 NOT NULL,
	user_id int4 NOT NULL,
	amount numeric(10, 2) NOT NULL,
	effective_from date NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	created_by varchar(255) NULL,
	CONSTRAINT user_salaries_pkey PRIMARY KEY (id),
	CONSTRAINT user_salaries_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);