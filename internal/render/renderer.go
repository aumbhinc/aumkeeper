package render

import (
	"html/template"
	"log"
	"net/http"
	 "strings"

)

type Renderer struct {
	Templates *template.Template
}

// RenderData is your unified contract for all pages
type RenderData struct {
	Layout interface{}
	Data   interface{}
}

// OK renders a page using centralized template resolution
func (r *Renderer) OK(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {

	// =========================================================
	// RESOLVE TEMPLATE NAME
	// Convention: page -> page_content
	// =========================================================
	resolvedPage := page + "_content"

	log.Println("=================================================")
	log.Println("🧭 RENDER START")
	log.Println("➡️ PAGE INPUT:", page)
	log.Println("➡️ TEMPLATE TARGET:", resolvedPage)
	log.Println("=================================================")

	// =========================================================
	// EXECUTE TEMPLATE INTO BUFFER
	// =========================================================
	var pageBuffer strings.Builder

	err := r.Templates.ExecuteTemplate(&pageBuffer, resolvedPage, data)
	if err != nil {
		log.Println("❌ TEMPLATE EXECUTION ERROR:", err)
		http.Error(w, "Page render error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("📦 PAGE BUFFER SIZE:", len(pageBuffer.String()))

	// =========================================================
	// FINAL OUTPUT TO RESPONSE
	// =========================================================
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(pageBuffer.String()))
}