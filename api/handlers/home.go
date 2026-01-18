package handlers

import (
	"aumkeeper/api/viewdata"
	"bytes"
	"html/template"
	"net/http"
	"time"
)

func HomeHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer

		// Render the "home_content" template into a buffer
		if err := t.ExecuteTemplate(&buf, "home_content", nil); err != nil {
			http.Error(w, "Error rendering page content: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare layout data
		data := viewdata.Layout{
			Title:       "Home",
			Description: "AumKeeper ERP platform for SMBs",
			Year:        time.Now().Year(),
			PageContent: template.HTML(buf.String()), // Safe raw HTML injection
		}

		// Render the base template with the page content
		if err := t.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Error rendering base template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
