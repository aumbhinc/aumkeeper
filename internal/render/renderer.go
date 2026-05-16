package render

import (
	"html/template"
	"log"
	"net/http"
)

type Renderer struct {
	tmpl *template.Template
}

// MUST MATCH handler Data contract safely
type RenderData struct {
	Title       string
	Description string
	Page        string
	Data        any
}

func NewRenderer(t *template.Template) *Renderer {
	return &Renderer{tmpl: t}
}

func (r *Renderer) OK(w http.ResponseWriter, page string, data *RenderData) {

	// ---------------------------------------------------------
	// CRITICAL FIX:
	// map logical page → actual template block name
	// ---------------------------------------------------------
	templateMap := map[string]string{
		"home":      "home_content",
		"dashboard": "dashboard_content",
		"staffs":    "staffs_content",
		"addstaff":  "addstaff_content",
	}

	tmplName, ok := templateMap[page]
	if !ok {
		log.Printf("PAGE TEMPLATE ERROR: unknown page mapping: %s", page)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	// ---------------------------------------------------------
	// EXECUTE SAFE TEMPLATE
	// ---------------------------------------------------------
	err := r.tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		log.Printf("PAGE TEMPLATE ERROR: %v", err)
		http.Error(w, "template render error", http.StatusInternalServerError)
		return
	}
}