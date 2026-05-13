package handlers

import (
	"html/template"
	"net/http"
	"time"

	"aumkeeper/api/templates"
	"aumkeeper/api/viewdata"
)

func StaffsHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		data := viewdata.Layout{
			Title:       "Staff Manager",
			Description: "Manage staff, roles, and payroll",
			Year:        time.Now().Year(),

			// 🔥 CORE ENGINE KEY
			Page: "staffs",

			IncludeStaffJS: true,
		}

		// 🔥 RESOLVE TEMPLATE VIA REGISTRY
		templates.ResolvePage(&data)

		// 🔥 SINGLE PASS RENDER
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Error rendering staff page: "+err.Error(), http.StatusInternalServerError)
		}
	})
}