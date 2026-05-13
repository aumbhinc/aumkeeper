package templates

import (
	"html/template"
	"net/http"
)

func RenderTemplate(t *template.Template, w http.ResponseWriter, name string, data any) error {
	return t.ExecuteTemplate(w, name, data)
}