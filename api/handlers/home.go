package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	"aumkeeper/api/viewdata"
)

func HomeHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Render home_content into buffer
		var buf bytes.Buffer
		if err := t.ExecuteTemplate(&buf, "home_content", nil); err != nil {
			http.Error(w, "Error rendering home content: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := viewdata.Layout{
			Title:       "Home",
			Description: "AumKeeper ERP platform for SMBs",
			Year:        time.Now().Year(),
			PageContent: template.HTML(buf.String()),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Error rendering home page: "+err.Error(), http.StatusInternalServerError)
		}
	})
}
