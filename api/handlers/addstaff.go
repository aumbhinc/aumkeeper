package handlers

import (
	"log"
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

func (h *AddStaffHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {

	case http.MethodGet:

		log.Println("🟢 GET /addstaff")

		h.Renderer.OK(w, "addstaff", &render.RenderData{
			Title:       "Add Staff",
			Description: "Create employee onboarding record",
			Page:        "addstaff",
			FormData:    render.FormData{},
			Errors:      map[string]string{},
		})
		return

	case http.MethodPost:

		log.Println("🔥 POST /addstaff HIT")

		if err := r.ParseForm(); err != nil {
			log.Println("❌ ParseForm error:", err)

			h.renderError(
				w,
				"Invalid form submission",
				domain.Staff{},
			)
			return
		}

		staff := domain.Staff{
			FirstName:  r.FormValue("firstName"),
			MiddleName: r.FormValue("middleName"),
			LastName:   r.FormValue("lastName"),
			Role:       r.FormValue("role"),

			Email:       r.FormValue("email"),
			PhoneNumber: r.FormValue("phoneNumber"),

			Street:  r.FormValue("street"),
			City:    r.FormValue("city"),
			State:   r.FormValue("state"),
			ZipCode: r.FormValue("zipCode"),

			SSN: r.FormValue("ssn"),

			DependentClaims: parseInt(
				r.FormValue("dependentClaims"),
			),

			Wage: parseFloat(
				r.FormValue("wage"),
			),

			PaymentFrequency: r.FormValue("paymentFrequency"),
			Comments:         r.FormValue("comments"),
		}

		log.Printf(
			"📩 Staff Submission: %s %s | Phone=%s | Role=%s",
			staff.FirstName,
			staff.LastName,
			staff.PhoneNumber,
			staff.Role,
		)

		created, err := h.StaffService.CreateStaff(
			r.Context(),
			staff,
		)

		if err != nil {
			log.Println(
				"❌ StaffService.CreateStaff error:",
				err,
			)

			h.renderError(
				w,
				err.Error(),
				staff,
			)
			return
		}

		log.Printf(
			"✅ Staff Created: %+v",
			created,
		)

		http.Redirect(
			w,
			r,
			"/staffs",
			http.StatusSeeOther,
		)
		return

	default:
		http.Error(
			w,
			"Method Not Allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *AddStaffHandler) renderError(
	w http.ResponseWriter,
	msg string,
	staff domain.Staff,
) {
	h.Renderer.OK(w, "addstaff", &render.RenderData{
		Title:       "Add Staff",
		Description: msg,
		Page:        "addstaff",

		FormData: render.FormData{
			FirstName:  staff.FirstName,
			MiddleName: staff.MiddleName,
			LastName:   staff.LastName,
			Role:       staff.Role,

			Email:       staff.Email,
			PhoneNumber: staff.PhoneNumber,

			Street:  staff.Street,
			City:    staff.City,
			State:   staff.State,
			ZipCode: staff.ZipCode,

			SSN: staff.SSN,

			DependentClaims: strconv.Itoa(
				staff.DependentClaims,
			),

			Wage: strconv.FormatFloat(
				staff.Wage,
				'f',
				2,
				64,
			),

			PaymentFrequency: staff.PaymentFrequency,
			Comments:         staff.Comments,
		},

		Errors: map[string]string{
			"form": msg,
		},
	})
}

func parseInt(s string) int {
	if s == "" {
		return 0
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return v
}

func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return v
}