package main

import (
    "html/template"
    "log"
    "net/http"
    "aumkeeper/api"
)

// Global template variable
var templates *template.Template

// LoadTemplates parses templates and registers helper functions like "dict"
func loadTemplates() {
    templates = template.Must(template.New("").Funcs(template.FuncMap{
        "dict": func(values ...interface{}) (map[string]interface{}, error) {
            if len(values)%2 != 0 {
                return nil, fmt.Errorf("dict expects even number of arguments")
            }
            m := make(map[string]interface{}, len(values)/2)
            for i := 0; i < len(values); i += 2 {
                key, ok := values[i].(string)
                if !ok {
                    return nil, fmt.Errorf("dict keys must be strings")
                }
                m[key] = values[i+1]
            }
            return m, nil
        },
    }).ParseGlob("api/templates/*.html"))
}

func main() {
    log.Println("🔄 Starting AumKeeper server...")

    // Load templates
    loadTemplates()
    api.Templates = templates // pass to your api package if needed

    // Register handlers
    api.RegisterRoutes()

    // Start server
    log.Println("✅ Server listening on http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("❌ Server failed: %v", err)
    }
}
