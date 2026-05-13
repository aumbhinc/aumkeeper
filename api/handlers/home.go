package handlers

import (
	"html/template"
	"net/http"
	"time"

	"aumkeeper/api/templates"
	"aumkeeper/api/viewdata"
)

func HomeHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		data := viewdata.Layout{
			Title:       "Home",
			Description: "AumKeeper ERP platform for SMBs",
			Year:        time.Now().Year(),

			// 🔥 CORE ROUTING KEY
			Page: "home",
		}

		// 🔥 RESOLVE TEMPLATE VIA REGISTRY
		templates.ResolvePage(&data)

		// 🔥 SINGLE PASS RENDER
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Error rendering home page: "+err.Error(), http.StatusInternalServerError)
		}
	})
}