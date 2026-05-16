package handlers

import (
	"net/http"
	"time"

	"aumkeeper/api/viewdata"
	"aumkeeper/internal/render"
)

type DashboardHandler struct {
	Renderer *render.Renderer
}

func NewDashboardHandler(r *render.Renderer) *DashboardHandler {
	return &DashboardHandler{Renderer: r}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	layout := viewdata.Layout{
		Title:       "Dashboard",
		Description: "AumKeeper Executive Control Panel",
		Year:        time.Now().Year(),
		Page:        "dashboard",
		IncludeDashboardJS: true,
	}

	h.Renderer.OK(w, "dashboard", &render.RenderData{
		Title:       layout.Title,
		Description: layout.Description,
		Page:        "dashboard",
		Data:        layout,
	})
}