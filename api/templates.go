package api

import "html/template"

// FuncMap exposes template helper functions
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"dict": func(values ...any) map[string]any {
			m := make(map[string]any)
			for i := 0; i < len(values); i += 2 {
				key, _ := values[i].(string)
				m[key] = values[i+1]
			}
			return m
		},
	}
}
