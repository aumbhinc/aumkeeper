package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	"aumkeeper/api/viewdata"
)

func StaffsHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Render staffs_content into buffer
		var buf bytes.Buffer
		if err := t.ExecuteTemplate(&buf, "staffs_content", nil); err != nil {
			http.Error(w, "Error rendering staffs content: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := viewdata.Layout{
			Title:          "Staff Manager",
			Description:    "Manage staff, roles, and payroll",
			IncludeStaffJS: true,
			Year:           time.Now().Year(),
			PageContent:    template.HTML(buf.String()),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Error rendering staff page: "+err.Error(), http.StatusInternalServerError)
		}
	})
}
