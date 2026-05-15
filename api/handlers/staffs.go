package handlers

import (
	"net/http"
	"time"

	"aumkeeper/api/templates"
	"aumkeeper/api/viewdata"
	"aumkeeper/internal/render"
)

type StaffsHandler struct {
	Renderer render.Renderer
}

func NewStaffsHandler(
	r render.Renderer,
) *StaffsHandler {

	return &StaffsHandler{
		Renderer: r,
	}
}

func (h *StaffsHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)

		return
	}

	data := viewdata.Layout{
		Title:       "Staff Manager",
		Description: "Manage staff, roles, and payroll",
		Year:        time.Now().Year(),

		IncludeStaffJS: true,
	}

	templates.ResolvePage(&data)

	h.Renderer.OK(
		w,
		"staffs",
		&render.RenderData{
			Title:       data.Title,
			Description: data.Description,
			Data:        &data,
		},
	)
}