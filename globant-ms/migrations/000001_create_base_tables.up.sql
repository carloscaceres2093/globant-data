-- Active: 1670342027352@@127.0.0.1@5432@postgres
CREATE TABLE IF NOT EXISTS jobs (
    id bigserial,
    job_id  bigint UNIQUE,
    job_name VARCHAR,
    created_at timestamptz,
    updated_at timestamptz,
    PRIMARY KEY (id));


CREATE INDEX IF NOT EXISTS job_name_index
    on jobs (job_name);


CREATE TABLE IF NOT EXISTS departments (
    id bigserial,
    department_id bigint UNIQUE,
    department_name VARCHAR,
    created_at timestamptz,
    updated_at timestamptz,
    PRIMARY KEY (id));

CREATE INDEX IF NOT EXISTS department_name_index
    on departments (department_name);


CREATE TABLE IF NOT EXISTS employees (
    id bigserial,
    employee_id int,
    employee_name VARCHAR,
    department_id bigint,
    job_id bigint,
    hire_time timestamptz,
    created_at timestamptz,
    updated_at timestamptz,
    PRIMARY KEY (id),
    CONSTRAINT fk_departments
    FOREIGN KEY (department_id)
    REFERENCES departments(department_id),
    CONSTRAINT fk_jobs
    FOREIGN KEY (job_id)
    REFERENCES jobs(job_id)
);

CREATE INDEX IF NOT EXISTS employee_id_index
    on employees (employee_id);