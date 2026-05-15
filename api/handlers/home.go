package handlers

import (
	"net/http"
	"time"

	"aumkeeper/api/templates"
	"aumkeeper/api/viewdata"
	"aumkeeper/internal/render"
)

type HomeHandler struct {
	Renderer *render.Renderer
}

func NewHomeHandler(
	renderer *render.Renderer,
) *HomeHandler {
	return &HomeHandler{
		Renderer: renderer,
	}
}

func (h *HomeHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	data := &viewdata.Layout{
		Title:       "Home",
		Description: "AumKeeper ERP platform for SMBs",
		Year:        time.Now().Year(),

		// CORE ROUTING KEY
		Page: "home",
	}

	// Resolve page metadata/template registry
	templates.ResolvePage(data)

	// SINGLE CENTRALIZED RENDER PIPELINE
	h.Renderer.OK(
		w,
		"home",
		&render.RenderData{
			Title:       data.Title,
			Description: data.Description,
			Page:        data.Page,
			Data:        data,
		},
	)
}