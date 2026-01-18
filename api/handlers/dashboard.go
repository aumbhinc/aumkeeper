package handlers

import (
	"html/template"
	"net/http"
	"time"
)

func DashboardHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "dashboard", map[string]any{
			"Title":              "Dashboard",
			"Year":               time.Now().Year(),
			"IncludeDashboardJS": true,
		})
	})
}
