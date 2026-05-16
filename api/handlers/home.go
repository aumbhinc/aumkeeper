package handlers

import (
	"net/http"
	"time"

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

	// =========================================================
	// STRICT GET ONLY
	// =========================================================
	if r.Method != http.MethodGet {
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	// =========================================================
	// PAGE LAYOUT
	// =========================================================
	layout := viewdata.Layout{
		Title:       "Home",
		Description: "AumKeeper ERP platform for SMBs",
		Year:        time.Now().Year(),
	}

	// =========================================================
	// RENDER PAGE (FIXED CONTRACT)
	// =========================================================
	h.Renderer.OK(
		w,
		"home",
		&render.RenderData{
			Title:       layout.Title,
			Description: layout.Description,
			Page:        "home",
			Data:        layout,
		},
	)
}