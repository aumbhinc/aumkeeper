package render

import (
	"html/template"
	"log"
	"net/http"
)

type Renderer struct {
	tmpl *template.Template
}

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

	// =========================================================
	// 1. PAGE → TEMPLATE NAME MAPPING (CRITICAL LAYER)
	// =========================================================
	templateMap := map[string]string{
		"home":      "home_content",
		"dashboard": "dashboard_content",
		"staffs":    "staffs_content",
		"addstaff":  "addstaff_content",
	}

	tmplName, ok := templateMap[page]
	if !ok {
		log.Printf("❌ RENDER ERROR: unknown page key: %s", page)
		http.Error(w, "template mapping not found", http.StatusInternalServerError)
		return
	}

	// =========================================================
	// 2. DEBUG: PRINT REQUESTED PAGE + RESOLVED TEMPLATE
	// =========================================================
	log.Printf("📦 RENDER REQUEST: page=%s → template=%s", page, tmplName)

	// =========================================================
	// 3. DEBUG: DUMP ALL LOADED TEMPLATES (CRITICAL TRACE)
	// =========================================================
	log.Println("===== TEMPLATE REGISTRY DUMP START =====")
	for _, t := range r.tmpl.Templates() {
		log.Printf("🧩 TEMPLATE LOADED: %q", t.Name())
	}
	log.Println("===== TEMPLATE REGISTRY DUMP END =====")

	// =========================================================
	// 4. EXECUTE TEMPLATE WITH FULL ERROR TRACE
	// =========================================================
	err := r.tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		log.Printf("❌ TEMPLATE EXECUTION ERROR:")
		log.Printf("   page      : %s", page)
		log.Printf("   template  : %s", tmplName)
		log.Printf("   error     : %v", err)

		http.Error(w, "template render failed", http.StatusInternalServerError)
		return
	}

	// =========================================================
	// 5. SUCCESS TRACE
	// =========================================================
	log.Printf("✅ TEMPLATE RENDER SUCCESS: %s", tmplName)
}