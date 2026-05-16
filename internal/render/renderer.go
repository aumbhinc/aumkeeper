package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"aumkeeper/api/viewdata"
)

type Renderer struct {
	Templates *template.Template
}

type RenderData struct {
	Layout *viewdata.Layout
	Page   string
	Data   any
}

func NewRenderer(t *template.Template) *Renderer {
	return &Renderer{
		Templates: t,
	}
}

func (r *Renderer) OK(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {
	r.Render(w, page, data)
}

func (r *Renderer) Render(
	w http.ResponseWriter,
	page string,
	data *RenderData,
) {

	/*
		STAGE 1:
		Render page template into buffer
	*/

	var pageBuffer bytes.Buffer

	err := r.Templates.ExecuteTemplate(
		&pageBuffer,
		page,
		data,
	)

	if err != nil {
		log.Println("PAGE TEMPLATE ERROR:", err)

		http.Error(
			w,
			"Page render error",
			http.StatusInternalServerError,
		)

		return
	}

	/*
		STAGE 2:
		Inject into layout
	*/

	pageHTML := template.HTML(pageBuffer.String())

	switch page {

	case "dashboard":
		data.Layout.DashboardContent = pageHTML

	case "staffs":
		data.Layout.StaffsContent = pageHTML

	case "addstaff":
		data.Layout.AddStaffContent = pageHTML

	default:
		log.Println("UNKNOWN PAGE:", page)
	}

	/*
		STAGE 3:
		Render base layout
	*/

	err = r.Templates.ExecuteTemplate(
		w,
		"base",
		data,
	)

	if err != nil {

		log.Println("BASE TEMPLATE ERROR:", err)

		http.Error(
			w,
			"Layout render error",
			http.StatusInternalServerError,
		)

		return
	}
}