package handlers

import (
	"net/http"
	"time"

	"aumkeeper/api/templates"
	"aumkeeper/api/viewdata"
	"aumkeeper/internal/render"
)

type DashboardHandler struct {
	Renderer render.Renderer
}

func NewDashboardHandler(
	r render.Renderer,
) *DashboardHandler {

	return &DashboardHandler{
		Renderer: r,
	}
}

func (h *DashboardHandler) ServeHTTP(
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

	layout := viewdata.Layout{
		Title:       "Dashboard",
		Description: "AumKeeper Executive Control Panel",
		Year:        time.Now().Year(),

		IncludeDashboardJS: true,

		Stats: []viewdata.DashboardStat{
			{
				Name:   "Net Worth",
				Value:  "$128,430",
				Change: "+12.4%",
				Icon:   "account_balance_wallet",
			},
			{
				Name:   "Total Assets",
				Value:  "$542,900",
				Change: "+8.2%",
				Icon:   "savings",
			},
			{
				Name:   "Liabilities",
				Value:  "$214,470",
				Change: "-3.1%",
				Icon:   "payments",
			},
			{
				Name:   "Ledger Exceptions",
				Value:  "3",
				Change: "-1 today",
				Icon:   "warning",
			},
		},
	}

	templates.ResolvePage(&layout)

	h.Renderer.OK(
		w,
		"dashboard",
		&render.RenderData{
			Title:       layout.Title,
			Description: layout.Description,
			Data:        &layout,
		},
	)
}