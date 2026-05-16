package handlers

import (
	"net/http"
	"strconv"
	"time"

	"aumkeeper/internal/domain"
	"aumkeeper/internal/render"
	"aumkeeper/internal/services"
)

type AddStaffForm struct {
	FirstName        string
	MiddleName       string
	LastName         string
	Role             string
	Email            string
	PhoneNumber      string
	Street           string
	City             string
	ZipCode          string
	SSN              string
	TaxFileStatus    string
	DependentClaims  int
	Wage             float64
	PaymentFrequency string
	Comments         string
}

type AddStaffHandler struct {
	Renderer     render.Renderer
	StaffService *services.StaffService
}

func NewAddStaffHandler(
	renderer render.Renderer,
	staffService *services.StaffService,
) *AddStaffHandler {
	return &AddStaffHandler{
		Renderer:     renderer,
		StaffService: staffService,
	}
}

func (h *AddStaffHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	// =========================================================
	// GET -> RENDER PAGE (IMPORTANT FIX HERE)
	// =========================================================
	case http.MethodGet:

		form := AddStaffForm{}

		h.Renderer.OK(w, "addstaff", &render.RenderData{
			Title:       "Add Staff",
			Description: "Create employee onboarding record for AumKeeper ERP",
			Page:        "addstaff",
			FormData:    form, // 🔥 CRITICAL FIX
			Errors:      map[string]string{},
		})

	// =========================================================
	// POST -> CREATE STAFF
	// =========================================================
	case http.MethodPost:

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dependentClaims, _ := strconv.Atoi(r.FormValue("dependentClaims"))
		wage, _ := strconv.ParseFloat(r.FormValue("wage"), 64)

		form := AddStaffForm{
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

		// =====================================================
		// SIMPLE VALIDATION (CAN EXPAND LATER)
		// =====================================================
		errors := map[string]string{}

		if form.FirstName == "" {
			errors["firstName"] = "First name required"
		}
		if form.Email == "" {
			errors["email"] = "Email required"
		}

		// =====================================================
		// IF ERRORS -> RE-RENDER FORM (CRITICAL FIX)
		// =====================================================
		if len(errors) > 0 {
			h.Renderer.OK(w, "addstaff", &render.RenderData{
				Title:       "Add Staff",
				Description: "Fix validation errors",
				Page:        "addstaff",
				FormData:    form,   // 🔥 KEEP USER INPUT
				Errors:      errors, // 🔥 SHOW ERRORS
			})
			return
		}

		// =====================================================
		// DOMAIN OBJECT (DB LAYER)
		// =====================================================
		staff := domain.Staff{
			FirstName:        form.FirstName,
			MiddleName:       form.MiddleName,
			LastName:         form.LastName,
			Role:             form.Role,
			Email:            form.Email,
			PhoneNumber:      form.PhoneNumber,
			Street:           form.Street,
			City:             form.City,
			ZipCode:          form.ZipCode,
			SSN:              form.SSN,
			TaxFileStatus:    form.TaxFileStatus,
			DependentClaims:  form.DependentClaims,
			Wage:             form.Wage,
			PaymentFrequency: form.PaymentFrequency,
			Comments:         form.Comments,
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
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}