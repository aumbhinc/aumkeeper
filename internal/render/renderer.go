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

// Stable renderer contract
type RenderData struct {
	Title       string
	Description string
	Page        string

	// Flexible page-specific payload
	Data map[string]any
}

func (r *Renderer) Render(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {
	// Renderer safety
	if r == nil || r.Tmpl == nil {
		http.Error(
			w,
			"renderer not initialized",
			http.StatusInternalServerError,
		)
		return
	}

	// Nil safety
	if data == nil {
		data = &RenderData{}
	}

	// Data bucket safety
	if data.Data == nil {
		data.Data = map[string]any{}
	}

	// Debug logs
	if r.Debug {
		log.Println("🧠 [RENDER START]")
		log.Println("➡️ Base template:", page)
		log.Println("➡️ Page:", data.Page)
		log.Println("➡️ Title:", data.Title)
	}

	// Ensure template exists
	t := r.Tmpl.Lookup(page)
	if t == nil {
		errMsg := fmt.Sprintf(
			"TEMPLATE NOT FOUND: %s",
			page,
		)

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

	// Final payload sent to template
	payload := map[string]any{
		"Title":       data.Title,
		"Description": data.Description,
		"Page":        data.Page,
		"Data":        data.Data,
	}

	// BUFFERED RENDERING
	// prevents partial writes + double WriteHeader panic
	var buf bytes.Buffer

	err := t.Execute(&buf, payload)
	if err != nil {
		log.Println("❌ TEMPLATE EXEC ERROR:", err)

		http.Error(
			w,
			"template execution error",
			http.StatusInternalServerError,
		)

		return
	}

	// Success
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

// OK shortcut
func (r *Renderer) OK(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {
	r.Render(w, page, data)
}

// Created shortcut
func (r *Renderer) Created(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {
	w.WriteHeader(http.StatusCreated)
	r.Render(w, page, data)
}

// NotFound shortcut
func (r *Renderer) NotFound(w http.ResponseWriter) {
	http.Error(
		w,
		"page not found",
		http.StatusNotFound,
	)
}

// ServerError shortcut
func (r *Renderer) ServerError(
	w http.ResponseWriter,
	err error,
) {
	log.Println("❌ SERVER ERROR:", err)

	http.Error(
		w,
		"internal server error",
		http.StatusInternalServerError,
	)
}