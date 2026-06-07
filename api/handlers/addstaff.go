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

	case http.MethodGet:
		h.Renderer.OK(w, "addstaff", &render.RenderData{
			Title:       "Add Staff",
			Description: "Create employee onboarding record",
			Page:        "addstaff",
			FormData:    render.FormData{},
			Errors:      map[string]string{},
		})
		return

	case http.MethodPost:

		if err := r.ParseForm(); err != nil {
			h.renderError(w, "Invalid form data", domain.Staff{})
			return
		}

		staff := domain.Staff{
			FirstName:  r.FormValue("firstName"),
			MiddleName: r.FormValue("middleName"),
			LastName:   r.FormValue("lastName"),

			Role: r.FormValue("role"),

			Email:       r.FormValue("email"),
			PhoneNumber: r.FormValue("phoneNumber"),

			Street:  r.FormValue("street"),
			City:    r.FormValue("city"),
			State:   r.FormValue("state"),
			ZipCode: r.FormValue("zipCode"),

			SSN: r.FormValue("ssn"),

			DependentClaims: parseInt(r.FormValue("dependentClaims")),
			Wage:            parseFloat(r.FormValue("wage")),

			PaymentFrequency: r.FormValue("paymentFrequency"),

			Comments: r.FormValue("comments"),
		}

		_, err := h.StaffService.CreateStaff(r.Context(), staff)
		if err != nil {
			h.renderError(w, err.Error(), staff)
			return
		}

		http.Redirect(w, r, "/staffs", http.StatusSeeOther)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *AddStaffHandler) renderError(
	w http.ResponseWriter,
	msg string,
	form domain.Staff,
) {
	h.Renderer.OK(w, "addstaff", &render.RenderData{
		Title:       "Add Staff",
		Description: msg,
		Page:        "addstaff",

		FormData: render.FormData{
			FirstName:  form.FirstName,
			MiddleName: form.MiddleName,
			LastName:   form.LastName,

			Role: form.Role,

			Email:       form.Email,
			PhoneNumber: form.PhoneNumber,

			Street:  form.Street,
			City:    form.City,
			State:   form.State,
			ZipCode: form.ZipCode,

			SSN: form.SSN,

			DependentClaims: strconv.Itoa(form.DependentClaims),
			Wage:            strconv.FormatFloat(form.Wage, 'f', 2, 64),

			PaymentFrequency: form.PaymentFrequency,

			Comments: form.Comments,
		},

		Errors: map[string]string{
			"form": msg,
		},
	})
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}