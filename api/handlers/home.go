package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

func HomeHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		data := map[string]any{
			"Title":       "Home",
			"Year":        time.Now().Year(),
			"PageContent": "home_content",
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := t.ExecuteTemplate(w, "base", data); err != nil {
			log.Println("home render error:", err)
			http.Error(w, "Template rendering error", http.StatusInternalServerError)
		}
	})
}
