package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Renderer struct {
	Tmpl *template.Template
	Debug bool
}

// RenderData = ONLY stable fields (NO Layout, NO FormData here)
type RenderData struct {
	Title       string
	Description string
	Page        string

	Data map[string]any // flexible payload bucket (IMPORTANT FIX)
}

// Render main entry
func (r *Renderer) Render(w http.ResponseWriter, page string, data *RenderData) {
	if r == nil || r.Tmpl == nil {
		http.Error(w, "renderer not initialized", http.StatusInternalServerError)
		return
	}

	if data == nil {
		data = &RenderData{}
	}

	if data.Data == nil {
		data.Data = map[string]any{}
	}

	// DEBUG LOGGING
	if r.Debug {
		log.Println("🧠 Render() called")
		log.Println("➡ page:", page)
		log.Println("➡ title:", data.Title)
	}

	tmplName := page

	// Safety: prevent undefined template panic
	t := r.Tmpl.Lookup(tmplName)
	if t == nil {
		errMsg := fmt.Sprintf("TEMPLATE NOT FOUND: %s", tmplName)
		log.Println("❌", errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	err := t.Execute(w, data)
	if err != nil {
		log.Println("❌ TEMPLATE EXEC ERROR:", err)
		http.Error(w, "template execution error", http.StatusInternalServerError)
		return
	}

	if r.Debug {
		log.Println("✅ Render success:", page)
	}
}

// Shortcut
func (r *Renderer) OK(w http.ResponseWriter, page string, data *RenderData) {
	r.Render(w, page, data)
}