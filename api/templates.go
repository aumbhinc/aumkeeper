package api

import (
	"fmt"
	"html/template"
	"log"
)

// FuncMap returns reusable template helpers for all templates.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		// dict helper: create a map for template usage
		// Usage: {{ $m := dict "key1" val1 "key2" val2 }}
		"dict": func(values ...any) map[string]any {
			m := make(map[string]any)
			if len(values)%2 != 0 {
				log.Println("⚠️ dict helper: uneven number of arguments, last key ignored")
			}
			for i := 0; i < len(values)-1; i += 2 {
				key, ok := values[i].(string)
				if !ok {
					log.Printf("⚠️ dict helper: key at index %d is not a string, skipping\n", i)
					continue
				}
				m[key] = values[i+1]
			}
			return m
		},

		// currency helper: format float64 as USD currency
		// Usage: {{ currency 1234.56 }} → $1,234.56
		"currency": func(v float64) string {
			return fmt.Sprintf("$%.2f", v)
		},

		// optional: safe default helper for template data
		// Usage: {{ default "fallback" .SomeValue }}
		"default": func(def any, val any) any {
			if val == nil {
				return def
			}
			return val
		},
	}
}
