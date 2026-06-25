CREATE TABLE staff (
    id BIGSERIAL PRIMARY KEY,

    -- Personal Information
    first_name TEXT NOT NULL,
    middle_name TEXT,
    last_name TEXT NOT NULL,
    job_role TEXT NOT NULL,

    email TEXT NOT NULL UNIQUE,
    phone_number TEXT,

    -- Address Information
    street_address TEXT,
    city TEXT,
    state TEXT,
    zip_code TEXT,

    -- Tax & Payroll Information
    ssn TEXT,
    dependents_claimed INTEGER DEFAULT 0,
    wage NUMERIC(12,2),
    payment_frequency TEXT,

    -- Internal Notes
    notes TEXT,

    -- Workforce Status
    status TEXT NOT NULL DEFAULT 'active',

    -- Audit Fields
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_staff_email ON staff(email);
CREATE INDEX idx_staff_last_name ON staff(last_name);
CREATE INDEX idx_staff_job_role ON staff(job_role);
CREATE INDEX idx_staff_status ON staff(status);