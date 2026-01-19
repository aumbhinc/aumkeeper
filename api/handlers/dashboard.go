package handlers

import (
	"aumkeeper/api/viewdata"
	"bytes"
	"html/template"
	"net/http"
	"time"
)

func DashboardHandler(t *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		page := viewdata.Dashboard{
			Stats: []viewdata.Stat{
				{"Net Worth", "$100MM", "+5%", "account_balance"},
				{"Total Assets", "$72MM", "+3%", "account_balance_wallet"},
				{"Total Liabilities", "$22MM", "-2%", "request_page"},
				{"Ledger Exceptions", "12", "-8%", "error_outline"},
			},
			RightPanelSections: map[string][]viewdata.PanelItem{
				"Accounts": {
					{"savings", "Asset Accounts", "success", "+39%"},
					{"payments", "Liability Accounts", "warn", "-9%"},
					{"monetization_on", "Revenue Accounts", "success", "+21%"},
					{"monetization_on", "Equity Accounts", "success", "+12%"},
					{"monetization_on", "Expense Accounts", "success", "+7%"},
					{"monetization_on", "Contra Accounts", "success", "+3%"},
				},
				"Operations": {
					{"add", "Pay Out", "warn", "2 pending"},
					{"add", "Add Product", "", ""},
				},
				"Sales Analytics": {
					{"shopping_cart", "Online Orders", "warn", "5 unfulfilled"},
					{"person", "New Customers", "success", "+12%"},
				},
			},
		}

		var buf bytes.Buffer
		t.ExecuteTemplate(&buf, "dashboard_content", page)

		layout := viewdata.Layout{
			Title:              "Dashboard",
			Description:        "AumKeeper Business Control Center",
			Year:               time.Now().Year(),
			IncludeDashboardJS: true,
			PageContent:        template.HTML(buf.String()),
		}

		t.ExecuteTemplate(w, "base", layout)
	})
}
