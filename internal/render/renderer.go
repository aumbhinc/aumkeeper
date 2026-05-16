package render

import (
	"html/template"
	"log"
	"net/http"
)

//
// =========================
// CORE DATA CONTRACT
// =========================
//

type RenderData struct {
	Title       string
	Description string
	Page        string

	// Primary payload (used by templates)
	Data any

	// 🔥 Strongly used by forms like AddStaff
	FormData any
	Errors   map[string]string

	// Feature flags
	IncludeDashboardJS bool
	IncludeStaffJS     bool

	// UI slots
	Sidebar any
	Topbar  any
	User    any
}

//
// =========================
// RENDER INTERFACE
// =========================
//

type Renderer interface {
	Render(w http.ResponseWriter, status int, page string, data *RenderData)
	OK(w http.ResponseWriter, page string, data *RenderData)
	Created(w http.ResponseWriter, page string, data *RenderData)
	NotFound(w http.ResponseWriter)
	ServerError(w http.ResponseWriter)
}

//
// =========================
// IMPLEMENTATION
// =========================
//

type HTMLRenderer struct {
	templates *template.Template
}

func NewRenderer(t *template.Template) Renderer {
	return &HTMLRenderer{
		templates: t,
	}
}

//
// =========================
// CORE SAFE RENDER ENGINE
// =========================
//

func (r *HTMLRenderer) Render(
	w http.ResponseWriter,
	status int,
	page string,
	data *RenderData,
) {

	if data == nil {
		data = &RenderData{}
	}

	data.Page = page

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	// 🔥 STRICT TEMPLATE NAMING RULE
	// All pages must follow: page + "_content"
	templateName := page + "_content"

	tmpl := r.templates.Lookup(templateName)
	if tmpl == nil {
		log.Printf("WARN: missing template: %s (fallback to base.html)", templateName)

		err := r.templates.ExecuteTemplate(w, "base.html", data)
		if err != nil {
			log.Printf("FATAL TEMPLATE ERROR (base fallback): %v", err)
			http.Error(w, "Template render failure", http.StatusInternalServerError)
		}
		return
	}

	// EXECUTE PAGE TEMPLATE SAFELY
	err := tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Printf("TEMPLATE EXEC ERROR [%s]: %v", templateName, err)
		http.Error(w, "Template render failure", http.StatusInternalServerError)
		return
	}
}

//
// =========================
// HELPERS
// =========================
//

func (r *HTMLRenderer) OK(w http.ResponseWriter, page string, data *RenderData) {
	r.Render(w, http.StatusOK, page, data)
}

func (r *HTMLRenderer) Created(w http.ResponseWriter, page string, data *RenderData) {
	r.Render(w, http.StatusCreated, page, data)
}

func (r *HTMLRenderer) NotFound(w http.ResponseWriter) {
	r.Render(w, http.StatusNotFound, "404", &RenderData{
		Title: "Page Not Found",
	})
}

func (r *HTMLRenderer) ServerError(w http.ResponseWriter) {
	r.Render(w, http.StatusInternalServerError, "500", &RenderData{
		Title: "Internal Server Error",
	})
}