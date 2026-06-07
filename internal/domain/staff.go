package domain

import "time"

type Staff struct {
ID           int64  `json:"id"`
EmployeeCode string `json:"employee_code"`

```
FirstName  string `json:"first_name"`
MiddleName string `json:"middle_name"`
LastName   string `json:"last_name"`

Role string `json:"role"`

Email       string `json:"email"`
PhoneNumber string `json:"phone_number"`

Street string `json:"street"`
City   string `json:"city"`
State  string `json:"state"`
ZipCode string `json:"zip_code"`

SSN string `json:"ssn"`

TaxFileStatus string `json:"tax_file_status"`

DependentClaims int `json:"dependent_claims"`

Wage             float64 `json:"wage"`
PaymentFrequency string  `json:"payment_frequency"`

Comments string `json:"comments"`

Status string `json:"status"`

CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
```

}
