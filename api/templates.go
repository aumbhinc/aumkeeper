package api

import (
	"fmt"
	"html/template"
)

// FuncMap provides reusable template helpers
func FuncMap() template.FuncMap {
	return template.FuncMap{
		// dict helper: {{ dict "key" value "key2" value2 }}
		"dict": func(values ...any) map[string]any {
			m := make(map[string]any)
			for i := 0; i < len(values); i += 2 {
				if i+1 < len(values) {
					key, _ := values[i].(string)
					m[key] = values[i+1]
				}
			}
			return m
		},

		// currency helper
		"currency": func(v float64) string {
			return fmt.Sprintf("$%.2f", v)
		},
	}
}
