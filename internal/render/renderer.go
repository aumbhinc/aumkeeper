package render

import (
	"html/template"
	"log"
	"net/http"
)

type RenderData struct {
	Title       string
	Description string
	Page        string
	Data        any

	IncludeDashboardJS bool
	IncludeStaffJS     bool

	Sidebar any
	Topbar  any
	User    any
}

// CONTRACT

type Renderer interface {
	Render(
		w http.ResponseWriter,
		status int,
		page string,
		data *RenderData,
	)

	OK(
		w http.ResponseWriter,
		page string,
		data *RenderData,
	)

	Created(
		w http.ResponseWriter,
		page string,
		data *RenderData,
	)

	NotFound(w http.ResponseWriter)

	ServerError(w http.ResponseWriter)
}

// IMPLEMENTATION

type HTMLRenderer struct {
	templates *template.Template
}

func NewRenderer(
	templates *template.Template,
) Renderer {
	return &HTMLRenderer{
		templates: templates,
	}
}

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

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.WriteHeader(status)

	err := r.templates.ExecuteTemplate(
		w,
		"base.html",
		data,
	)

	if err != nil {

		log.Printf(
			"render error [%s]: %v",
			page,
			err,
		)

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}
}

func (r *HTMLRenderer) OK(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {
	r.Render(w, http.StatusOK, page, data)
}

func (r *HTMLRenderer) Created(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {
	r.Render(w, http.StatusCreated, page, data)
}

func (r *HTMLRenderer) NotFound(
	w http.ResponseWriter,
) {
	r.Render(
		w,
		http.StatusNotFound,
		"404",
		&RenderData{
			Title: "Page Not Found",
		},
	)
}

func (r *HTMLRenderer) ServerError(
	w http.ResponseWriter,
) {
	r.Render(
		w,
		http.StatusInternalServerError,
		"500",
		&RenderData{
			Title: "Internal Server Error",
		},
	)
}