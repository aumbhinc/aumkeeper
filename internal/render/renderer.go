package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type FormData struct {
	FirstName       string
	MiddleName      string
	LastName        string
	Role            string
	Email           string
	PhoneNumber     string
	Street          string
	City            string
	State           string
	ZipCode         string
	SSN             string
	DependentClaims string
	Wage            string
	PaymentFrequency string
	Comments        string
}

type RenderData struct {
	Title       string
	Description string
	Page        string

	FormData FormData
	Errors   map[string]string
}

type Renderer struct {
	Tmpl  *template.Template
	Debug bool
}

func NewRenderer(t *template.Template, debug bool) *Renderer {
	return &Renderer{
		Tmpl:  t,
		Debug: debug,
	}
}

func (r *Renderer) Render(w http.ResponseWriter, page string, data *RenderData) {
	if r == nil || r.Tmpl == nil {
		http.Error(w, "renderer not initialized", http.StatusInternalServerError)
		return
	}

	if data == nil {
		data = &RenderData{}
	}

	if data.Errors == nil {
		data.Errors = map[string]string{}
	}

	if r.Debug {
		log.Println("RENDER START:", page)
	}

	t := r.Tmpl.Lookup(page)
	if t == nil {
		http.Error(w, fmt.Sprintf("template not found: %s", page), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer

	if err := t.Execute(&buf, data); err != nil {
		http.Error(w, "template execution error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (r *Renderer) OK(w http.ResponseWriter, page string, data *RenderData) {
	r.Render(w, page, data)
}