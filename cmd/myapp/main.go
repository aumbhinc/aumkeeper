package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aumkeeper/api"
	"aumkeeper/api/handlers"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("🔄 Starting AumKeeper server...")

	// Load environment variables
	_ = godotenv.Load()

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleSignals(cancel)

	// Load templates with FuncMap
	templates := template.New("").Funcs(api.FuncMap())
	if _, err := templates.ParseGlob("api/templates/*.html"); err != nil {
		log.Fatalf("❌ Template parse error: %v", err)
	}
	log.Println("✅ Templates loaded")

	// HTTP multiplexer
	mux := http.NewServeMux()

	// Static assets
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("api/static")))
	mux.Handle("/static/", cacheControlMiddleware(static))

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Home page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, templates, "home", map[string]any{
			"Title": "Home",
			"Year":  time.Now().Year(),
		})
	})

	// Dashboard page
	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, templates, "dashboard", map[string]any{
			"Title":              "Dashboard",
			"Year":               time.Now().Year(),
			"IncludeDashboardJS": true,
		})
	})

	// Staffs page (dynamic handler)
	mux.Handle("/staffs", handlers.StaffsHandler(templates))

	// Server configuration
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

	// Wait for shutdown signal
	<-ctx.Done()
	gracefulShutdown(server, 10*time.Second)
}

// renderTemplate executes template safely
func renderTemplate(w http.ResponseWriter, tmpl *template.Template, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Println("❌ Render error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// getPort returns env port or default 8080
func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

// handleSignals cancels context on SIGINT/SIGTERM
func handleSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	cancel()
}

// gracefulShutdown shuts down the server with timeout
func gracefulShutdown(server *http.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("❌ Graceful shutdown error:", err)
	} else {
		log.Println("✅ Server gracefully stopped")
	}
}

// cacheControlMiddleware sets caching headers for static files
func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		next.ServeHTTP(w, r)
	})
}
