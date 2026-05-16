package handlers

import (
	"net/http"
	"time"

	"aumkeeper/api/viewdata"
	"aumkeeper/internal/render"
)

type StaffsHandler struct {
	Renderer *render.Renderer
}

func NewStaffsHandler(
	r *render.Renderer,
) *StaffsHandler {

	return &StaffsHandler{
		Renderer: r,
	}
}

func (h *StaffsHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {

	// =========================================================
	// STRICT GET ONLY
	// =========================================================
	if r.Method != http.MethodGet {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	// =========================================================
	// PAGE LAYOUT
	// =========================================================
	layout := viewdata.Layout{
		Title:       "Staff Manager",
		Description: "Manage staff, roles, and payroll",
		Year:        time.Now().Year(),

		Page: "staffs",

		IncludeStaffJS: true,
	}

	// =========================================================
	// RENDER PAGE (FIXED CONTRACT)
	// =========================================================
	h.Renderer.OK(
		w,
		"staffs",
		&render.RenderData{
			Title:       layout.Title,
			Description: layout.Description,
			Page:        layout.Page,

			// IMPORTANT FIX:
			// pass VALUE, not pointer (prevents structural drift)
			Data: layout,
		},
	)
}