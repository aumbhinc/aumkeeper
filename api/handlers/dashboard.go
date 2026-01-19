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
		var buf bytes.Buffer
		t.ExecuteTemplate(&buf, "dashboard_content", nil)

		data := viewdata.Layout{
			Title:       "Dashboard",
			Description: "AumKeeper ERP Dashboard",
			Year:        time.Now().Year(),
			PageContent: template.HTML(buf.String()),
			Stats: []viewdata.StatCard{
				{"Net Worth", "account_balance", "$100MM", "+5%"},
				{"Total Assets", "account_balance_wallet", "$100MM", "+5%"},
				{"Total Liabilities", "request_page", "$100MM", "+5%"},
				{"Ledger Exceptions", "error_outline", "$100MM", "+5%"},
			},
			RightPanelSections: map[string][]viewdata.FuncButton{
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

		t.ExecuteTemplate(w, "base", data)
	})
}
