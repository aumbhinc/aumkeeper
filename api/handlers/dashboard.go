package handlers

import (
	"html/template"
	"net/http"
	"time"
)

type DashboardHandler struct {
	Tmpl *template.Template
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Later: auth check middleware should guard this route

	data := map[string]any{
		"Title":              "Dashboard",
		"IncludeDashboardJS": true, // CRITICAL
		"Year":               time.Now().Year(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := h.Tmpl.ExecuteTemplate(w, "dashboard", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
