package render

import (
	"bytes"
	"html/template"
	"net/http"
	"sync"
)

type Renderer struct {
	tmpl *template.Template
	mu   sync.RWMutex
}

func NewRenderer(t *template.Template) *Renderer {
	return &Renderer{tmpl: t}
}

// =========================
// SINGLE SOURCE OF TRUTH
// =========================
type RenderData struct {
	// layout metadata (WAS MISSING → causing your error)
	Title       string
	Description string

	// page routing key for base.html
	Page string

	// generic payload
	Data any
}

func (r *Renderer) OK(w http.ResponseWriter, page string, data *RenderData) {
	r.render(w, http.StatusOK, page, data)
}

func (r *Renderer) Render(w http.ResponseWriter, status int, page string, data *RenderData) {
	r.render(w, status, page, data)
}

func (r *Renderer) ServerError(w http.ResponseWriter, page string, data *RenderData) {
	r.render(w, http.StatusInternalServerError, page, data)
}

func (r *Renderer) render(w http.ResponseWriter, status int, page string, data *RenderData) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	if data == nil {
		data = &RenderData{}
	}

	data.Page = page

	var buf bytes.Buffer

	err := r.tmpl.ExecuteTemplate(&buf, "base.html", data)
	if err != nil {
		http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}