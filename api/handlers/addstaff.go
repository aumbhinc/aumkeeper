package handlers

import (
	"net/http"
	"strconv"

	"aumkeeper/internal/domain"
	"aumkeeper/internal/render"
	"aumkeeper/internal/services"
)

type AddStaffHandler struct {
	Renderer     *render.Renderer
	StaffService *services.StaffService
}

func NewAddStaffHandler(
	renderer *render.Renderer,
	staffService *services.StaffService,
) *AddStaffHandler {
	return &AddStaffHandler{
		Renderer:     renderer,
		StaffService: staffService,
	}
}

func (h *AddStaffHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	// =========================
	// GET (SAFE EMPTY FORM)
	// =========================
	case http.MethodGet:

		h.Renderer.OK(w, "addstaff", &render.RenderData{
			Title:       "Add Staff",
			Description: "Create employee onboarding record",
			Page:        "addstaff_content",

			// SAFE DEFAULT FORM (NEVER NIL)
			Data: map[string]any{
				"FormData": domain.Staff{},
			},
		})

	// =========================
	// POST
	// =========================
	case http.MethodPost:

		if err := r.ParseForm(); err != nil {
			h.renderError(w, "Invalid form data", domain.Staff{})
			return
		}

		staff := domain.Staff{
			FirstName:        r.FormValue("firstName"),
			MiddleName:       r.FormValue("middleName"),
			LastName:         r.FormValue("lastName"),
			Role:              r.FormValue("role"),
			Email:             r.FormValue("email"),
			PhoneNumber:      r.FormValue("phoneNumber"),
			Street:            r.FormValue("street"),
			City:              r.FormValue("city"),
			ZipCode:           r.FormValue("zipCode"),
			SSN:               r.FormValue("ssn"),
			TaxFileStatus:     r.FormValue("taxFileStatus"),
			DependentClaims:   parseInt(r.FormValue("dependentClaims")),
			Wage:              parseFloat(r.FormValue("wage")),
			PaymentFrequency:  r.FormValue("paymentFrequency"),
			Comments:          r.FormValue("comments"),
		}

		_, err := h.StaffService.CreateStaff(r.Context(), staff)
		if err != nil {
			h.renderError(w, "Failed to create staff", staff)
			return
		}

		http.Redirect(w, r, "/staffs", http.StatusSeeOther)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

//
// =========================
// ERROR RENDER SAFE
// =========================
//

func (h *AddStaffHandler) renderError(w http.ResponseWriter, msg string, form domain.Staff) {

	h.Renderer.OK(w, "addstaff", &render.RenderData{
		Title:       "Add Staff",
		Description: msg,
		Page:        "addstaff_content",

		// PRESERVE FORM STATE (NO RESET ON ERROR)
		Data: map[string]any{
			"FormData": form,
		},
	})
}

//
// =========================
// PARSERS
// =========================
//

func parseInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func parseFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}