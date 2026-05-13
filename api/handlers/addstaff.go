package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"aumkeeper/api/templates"
	"aumkeeper/api/viewdata"
	"aumkeeper/internal/domain"
	"aumkeeper/internal/services"
)

type AddStaffHandler struct {
	Templates    *template.Template
	StaffService *services.StaffService
}

func NewAddStaffHandler(
	templates *template.Template,
	staffService *services.StaffService,
) *AddStaffHandler {
	return &AddStaffHandler{
		Templates:    templates,
		StaffService: staffService,
	}
}

func (h *AddStaffHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	// =========================================================
	// GET -> Render Add Staff Page
	// =========================================================
	case http.MethodGet:

		data := viewdata.Layout{
			Title:       "Add Staff",
			Description: "Create employee onboarding record for AumKeeper ERP",
			Year:        time.Now().Year(),

			// 🔥 CORE ENGINE KEY
			Page: "addstaff",

			IncludeStaffJS: true,
		}

		// 🔥 RESOLVE TEMPLATE
		templates.ResolvePage(&data)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := h.Templates.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	// =========================================================
	// POST -> Create Staff
	// =========================================================
	case http.MethodPost:

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dependentClaims, _ := strconv.Atoi(r.FormValue("dependentClaims"))
		wage, _ := strconv.ParseFloat(r.FormValue("wage"), 64)

		staff := domain.Staff{
			FirstName:        r.FormValue("firstName"),
			MiddleName:       r.FormValue("middleName"),
			LastName:         r.FormValue("lastName"),
			Role:             r.FormValue("role"),
			Email:            r.FormValue("email"),
			PhoneNumber:      r.FormValue("phoneNumber"),
			Street:           r.FormValue("street"),
			City:             r.FormValue("city"),
			ZipCode:          r.FormValue("zipCode"),
			SSN:              r.FormValue("ssn"),
			TaxFileStatus:    r.FormValue("taxFileStatus"),
			DependentClaims:  dependentClaims,
			Wage:             wage,
			PaymentFrequency: r.FormValue("paymentFrequency"),
			Comments:         r.FormValue("comments"),
		}

		_, err := h.StaffService.CreateStaff(r.Context(), staff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/staffs", http.StatusSeeOther)

	// =========================================================
	// METHOD NOT ALLOWED
	// =========================================================
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}