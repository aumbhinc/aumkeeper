package handlers

import (
	"net/http"
	"strconv"
	"time"

	"aumkeeper/api/viewdata"
	"aumkeeper/internal/domain"
	"aumkeeper/internal/render"
	"aumkeeper/internal/services"
)

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

func (h *AddStaffHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	// =========================================================
	// GET -> RENDER ADD STAFF PAGE
	// =========================================================
	case http.MethodGet:

		layout := viewdata.Layout{
			Title:       "Add Staff",
			Description: "Create employee onboarding record for AumKeeper ERP",
			Year:        time.Now().Year(),

			Page: "addstaff",

			IncludeStaffJS: true,
		}

		h.Renderer.OK(
			w,
			"addstaff",
			&render.RenderData{
				Title:       layout.Title,
				Description: layout.Description,
				Page:        layout.Page,
				Data:        &layout,
			},
		)

	// =========================================================
	// POST -> CREATE STAFF
	// =========================================================
	case http.MethodPost:

		if err := r.ParseForm(); err != nil {

			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)

			return
		}

		dependentClaims, err := strconv.Atoi(
			r.FormValue("dependentClaims"),
		)

		if err != nil {
			dependentClaims = 0
		}

		wage, err := strconv.ParseFloat(
			r.FormValue("wage"),
			64,
		)

		if err != nil {
			wage = 0
		}

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

		_, err = h.StaffService.CreateStaff(
			r.Context(),
			staff,
		)

		if err != nil {

			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)

			return
		}

		http.Redirect(
			w,
			r,
			"/staffs",
			http.StatusSeeOther,
		)

	// =========================================================
	// METHOD NOT ALLOWED
	// =========================================================
	default:

		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}