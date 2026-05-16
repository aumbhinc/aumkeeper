package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type RenderData struct {
	Title       string
	Description string
	Page        string

	Data  any
	Debug bool
}

type Renderer struct {
	Tmpl *template.Template
}

func NewRenderer(tmpl *template.Template) *Renderer {
	return &Renderer{
		Tmpl: tmpl,
	}
}

// =========================================================
// MAIN ENTRY
// =========================================================
func (r *Renderer) OK(w http.ResponseWriter, base string, data *RenderData) {

	log.Println("🧠 [RENDER START]")
	log.Println("➡️ Base template:", base)
	log.Println("➡️ Page:", data.Page)
	log.Println("➡️ Title:", data.Title)

	if data == nil {
		log.Println("❌ RenderData is NIL")
		http.Error(w, "render data missing", http.StatusInternalServerError)
		return
	}

	if data.Page == "" {
		log.Println("❌ PAGE IS EMPTY → template will fail")
		http.Error(w, "page not defined", http.StatusInternalServerError)
		return
	}

	// DEBUG: template existence check
	t := r.Tmpl.Lookup(base)
	if t == nil {
		log.Println("❌ BASE TEMPLATE NOT FOUND:", base)
		http.Error(w, "base template missing", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Base template found")

	// Execute
	err := t.Execute(w, data)
	if err != nil {
		log.Println("❌ TEMPLATE EXEC ERROR:", err)
		http.Error(w, fmt.Sprintf("template error: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("✅ [RENDER SUCCESS]")
}

// =========================================================
// HELPERS
// =========================================================
func (r *Renderer) NotFound(w http.ResponseWriter) {
	http.Error(w, "404 not found", http.StatusNotFound)
}

func (r *Renderer) ServerError(w http.ResponseWriter, err error) {
	log.Println("❌ SERVER ERROR:", err)
	http.Error(w, "server error", http.StatusInternalServerError)
}