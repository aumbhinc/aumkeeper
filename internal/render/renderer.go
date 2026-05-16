package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Renderer struct {
	Tmpl  *template.Template
	Debug bool
}

type RenderData struct {
	Title       string
	Description string
	Page        string

	Data map[string]any
}

func (r *Renderer) Render(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {
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

	if r.Debug {
		log.Println("🧠 [RENDER START]")
		log.Println("➡️ Base template:", page)
		log.Println("➡️ Page:", data.Page)
		log.Println("➡️ Title:", data.Title)
	}

	t := r.Tmpl.Lookup(page)
	if t == nil {
		errMsg := fmt.Sprintf("template not found: %s", page)
		log.Println("❌", errMsg)

		http.Error(
			w,
			errMsg,
			http.StatusInternalServerError,
		)
		return
	}

	if r.Debug {
		log.Println("✅ Base template found")
	}

	var buf bytes.Buffer

	err := t.Execute(&buf, data)
	if err != nil {
		log.Println("❌ TEMPLATE EXEC ERROR:", err)

		http.Error(
			w,
			"template execution error",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("❌ RESPONSE WRITE ERROR:", err)
		return
	}

	if r.Debug {
		log.Println("✅ [RENDER SUCCESS]")
	}
}

func (r *Renderer) OK(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {
	r.Render(w, page, data)
}