package api

import (
	"fmt"
	"html/template"
	"strings"
)

// FuncMap provides reusable template functions
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"currency": func(amount float64) string {
			return "$" + fmt.Sprintf("%.2f", amount)
		},
	}
}
