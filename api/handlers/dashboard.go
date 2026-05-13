package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	"aumkeeper/api/viewdata"
)

func DashboardHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		layout := viewdata.Layout{
			Title:              "Dashboard",
			Description:        "AumKeeper Executive Control Panel",
			Year:               time.Now().Year(),
			IncludeDashboardJS: true,

			Stats: []viewdata.DashboardStat{
				{Name: "Net Worth", Value: "$128,430", Change: "+12.4%", Icon: "account_balance_wallet"},
				{Name: "Total Assets", Value: "$542,900", Change: "+8.2%", Icon: "savings"},
				{Name: "Liabilities", Value: "$214,470", Change: "-3.1%", Icon: "payments"},
				{Name: "Ledger Exceptions", Value: "3", Change: "-1 today", Icon: "warning"},
			},

			RightPanelSections: map[string][]viewdata.RightPanelItem{
				"Accounts": {
					{Icon: "savings", Label: "Asset Accounts", Alert: "success", Value: "+39%"},
					{Icon: "payments", Label: "Liability Accounts", Alert: "warn", Value: "-9%"},
					{Icon: "monetization_on", Label: "Revenue Accounts", Alert: "success", Value: "+21%"},
					{Icon: "monetization_on", Label: "Equity Accounts", Alert: "success", Value: "+12%"},
					{Icon: "monetization_on", Label: "Expense Accounts", Alert: "success", Value: "+7%"},
					{Icon: "monetization_on", Label: "Contra Accounts", Alert: "success", Value: "+3%"},
				},
				"Operations": {
					{Icon: "add", Label: "Pay Out", Alert: "warn", Value: "2 pending"},
					{Icon: "add", Label: "Add Product"},
					{Icon: "add", Label: "Add Staff"},
					{Icon: "add", Label: "Add Customer"},
					{Icon: "add", Label: "Add Market"},
				},
				"Sales Analytics": {
					{Icon: "shopping_cart", Label: "Online Orders", Alert: "warn", Value: "5 unfulfilled"},
					{Icon: "person", Label: "New Customers", Alert: "success", Value: "+12%"},
					{Icon: "trending_up", Label: "Revenue Trends"},
					{Icon: "bar_chart", Label: "Sales Overview"},
				},
			},
		}

		// Render dashboard_content into buffer
		var buf bytes.Buffer
		if err := t.ExecuteTemplate(&buf, "dashboard_content", layout); err != nil {
			http.Error(w, "Error rendering dashboard content: "+err.Error(), http.StatusInternalServerError)
			return
		}

		layout.PageContent = template.HTML(buf.String())

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(w, "base", layout); err != nil {
			http.Error(w, "Error rendering base template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
