package handlers

import (
	"aumkeeper/api/viewdata"
	"bytes"
	"html/template"
	"net/http"
	"time"
)

// DashboardHandler renders the dashboard page
func DashboardHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Render dashboard_content into a buffer first
		var buf bytes.Buffer
		if err := t.ExecuteTemplate(&buf, "dashboard_content", nil); err != nil {
			http.Error(w, "Error rendering dashboard content: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare layout data for base template
		data := viewdata.Layout{
			Title:          "Dashboard",
			Description:    "AumKeeper Dashboard - View business insights",
			Year:           time.Now().Year(),
			PageContent:    template.HTML(buf.String()), // safe HTML
			IncludeDashboardJS: true,                     // loads dashboard CSS & JS
		}

		// Render the base template with the dashboard content
		if err := t.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Error rendering dashboard page: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
