package render

import (
	"html/template"
	"log"
	"net/http"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer(templates *template.Template) *Renderer {
	return &Renderer{
		templates: templates,
	}
}

// RenderData is the universal render payload
// that flows through the entire ERP UI system.
type RenderData struct {
	Title       string
	Description string

	// Page identifier from PageRegistry
	Page string

	// Actual page/module data
	Data any

	// Global UI flags
	IncludeDashboardJS bool
	IncludeStaffJS     bool

	// Future-safe extension fields
	Sidebar any
	Topbar  any
	User    any
}

// Render is the ONLY approved rendering entrypoint.
//
// ALL handlers MUST flow through this method.
// NO handler should directly execute templates.
func (r *Renderer) Render(
	w http.ResponseWriter,
	status int,
	page string,
	data *RenderData,
) {
	if data == nil {
		data = &RenderData{}
	}

	// Inject page automatically
	data.Page = page

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	err := r.templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {

		log.Printf("render error [%s]: %v", page, err)

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}
}

// Convenience helpers

func (r *Renderer) OK(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {
	r.Render(w, http.StatusOK, page, data)
}

func (r *Renderer) Created(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {
	r.Render(w, http.StatusCreated, page, data)
}

func (r *Renderer) NotFound(w http.ResponseWriter) {
	r.Render(w, http.StatusNotFound, "404", &RenderData{
		Title: "Page Not Found",
	})
}

func (r *Renderer) ServerError(w http.ResponseWriter) {
	r.Render(w, http.StatusInternalServerError, "500", &RenderData{
		Title: "Internal Server Error",
	})
}