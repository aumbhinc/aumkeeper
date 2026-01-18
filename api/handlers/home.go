package handlers

import (
	"aumkeeper/api/viewdata"
	"html/template"
	"net/http"
	"time"
)

func HomeHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var content template.HTML
		t.ExecuteTemplate(&content, "home_content", nil)

		data := viewdata.Layout{
			Title:       "Home",
			Description: "AumKeeper ERP platform for SMBs",
			Year:        time.Now().Year(),
			PageContent: content,
		}

		t.ExecuteTemplate(w, "base", data)
	})
}
