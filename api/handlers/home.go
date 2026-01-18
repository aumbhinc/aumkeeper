package handlers

import (
	"html/template"
	"net/http"
	"time"
)

func HomeHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "home", map[string]any{
			"Title": "Home",
			"Year":  time.Now().Year(),
		})
	})
}
