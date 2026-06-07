-- =========================================
-- AUMKEEPER INITIAL STAFF TABLE
-- =========================================

CREATE TABLE IF NOT EXISTS staff (
id BIGSERIAL PRIMARY KEY,

```
employee_code TEXT NOT NULL,

first_name TEXT NOT NULL,
middle_name TEXT,
last_name TEXT NOT NULL,

role TEXT,

email TEXT,
phone_number TEXT,

street TEXT,
city TEXT,
state TEXT,
zip_code TEXT,

ssn TEXT,
tax_file_status TEXT,

dependent_claims INTEGER DEFAULT 0,

wage NUMERIC DEFAULT 0,
payment_frequency TEXT,

comments TEXT,

status TEXT DEFAULT 'active',

created_at TIMESTAMP DEFAULT NOW(),
updated_at TIMESTAMP DEFAULT NOW()
```

);
