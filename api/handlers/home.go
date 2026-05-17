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

func NewHomeHandler(renderer *render.Renderer) *HomeHandler {
	return &HomeHandler{Renderer: renderer}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	layout := viewdata.Layout{
		Title:       "Home",
		Description: "AumKeeper ERP platform",
		Year:        time.Now().Year(),
		Page:        "home",
	}

	h.Renderer.OK(w, "home", &render.RenderData{
		Title:       layout.Title,
		Description: layout.Description,
		Page:        "home",
	})
}