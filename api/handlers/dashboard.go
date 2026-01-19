package handlers

import (
	"aumkeeper/api/viewdata"
	"bytes"
	"html/template"
	"net/http"
)

func DashboardHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Build Dashboard ViewData
		data := viewdata.Dashboard{
			Stats: []viewdata.DashboardStat{
				{Name: "Net Worth", Value: "$128,430", Change: "+12.4%", Icon: "account_balance_wallet"},
				{Name: "Total Assets", Value: "$542,900", Change: "+8.2%", Icon: "savings"},
				{Name: "Liabilities", Value: "$214,470", Change: "-3.1%", Icon: "payments"},
				{Name: "Ledger Exceptions", Value: "3", Change: "-1 today", Icon: "warning"},
			},

			RightPanelSections: map[string][]viewdata.DashboardPanelItem{
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
				},

				"Sales Analytics": {
					{Icon: "shopping_cart", Label: "Online Orders", Alert: "warn", Value: "5 unfulfilled"},
					{Icon: "person", Label: "New Customers", Alert: "success", Value: "+12%"},
				},
			},
		}

		// Render dashboard content into buffer
		var buf bytes.Buffer
		err := t.ExecuteTemplate(&buf, "dashboard_content", data)
		if err != nil {
			http.Error(w, "Error rendering dashboard content: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Inject into base layout
		layout := viewdata.Layout{
			Title:       "Dashboard",
			Description: "AumKeeper Executive Control Panel",
			PageContent: template.HTML(buf.String()),
		}

		err = t.ExecuteTemplate(w, "base", layout)
		if err != nil {
			http.Error(w, "Error rendering base template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
