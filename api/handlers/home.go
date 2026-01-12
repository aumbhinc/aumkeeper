package handlers

import (
	"html/template"
	"net/http"
	"time"
)

type HomeHandler struct {
	Tmpl *template.Template
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":       "Modern ERP for Growing Businesses",
		"Description": "AumKeeper is a modern all-in-one SaaS ERP platform for SMBs.",
		"Year":        time.Now().Year(),
		// IMPORTANT: Do NOT set IncludeDashboardJS
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := h.Tmpl.ExecuteTemplate(w, "home", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
