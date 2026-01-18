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
		t.ExecuteTemplate(&buf, "home_content", nil)

		data := viewdata.Layout{
			Title:       "Home",
			Description: "AumKeeper ERP platform for SMBs",
			Year:        time.Now().Year(),
			PageContent: template.HTML(buf.String()),
		}

		t.ExecuteTemplate(w, "base", data)
	})
}
