package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aumkeeper/api"
	"aumkeeper/api/handlers"

	"github.com/joho/godotenv"
	"html/template"
)

func main() {
	log.Println("🔄 Starting AumKeeper server...")
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleSignals(cancel)

	// Load templates with FuncMap
	templates := template.New("").Funcs(api.FuncMap())
	if _, err := templates.ParseGlob("api/templates/*.html"); err != nil {
		log.Fatalf("❌ Template parse error: %v", err)
	}
	log.Println("✅ Templates loaded")

	mux := http.NewServeMux()

	// Static files
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("api/static")))
	mux.Handle("/static/", cacheControlMiddleware(static))

	// Routes
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, templates, "home", map[string]any{
			"Title": "Home",
			"Year":  time.Now().Year(),
		})
	})
	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, templates, "dashboard", map[string]any{
			"Title":              "Dashboard",
			"Year":               time.Now().Year(),
			"IncludeDashboardJS": true,
		})
	})
	mux.Handle("/staffs", handlers.StaffsHandler(templates)) // dynamic staff page

	server := &http.Server{
		Addr:              ":" + getPort(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("🌐 Listening on :%s", getPort())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	gracefulShutdown(server, 10*time.Second)
}

// Helpers

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Println("❌ Render error:", err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

func handleSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	cancel()
}

func gracefulShutdown(server *http.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_ = server.Shutdown(ctx)
}

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		next.ServeHTTP(w, r)
	})
}
