package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// templates holds all parsed HTML templates
var templates *template.Template

func main() {
	fmt.Println("🔄 Starting AumKeeper server...")

	// Load .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found or failed to load")
	}

	// Graceful shutdown context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// OS signal handling for shutdown
	go handleSignals(cancel)

	// Load templates once
	loadTemplates("api/templates/*.html")

	// Setup router
	mux := http.NewServeMux()
	setupRoutes(mux)

	// Determine port
	port := getPort()

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("🌐 Server listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	gracefulShutdown(server, 10*time.Second)
}

/* -------------------- Helper Functions -------------------- */

func handleSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	log.Printf("🛑 Received signal: %v", sig)
	cancel()
}

func loadTemplates(pattern string) {
	var err error
	templates, err = template.ParseGlob(pattern)
	if err != nil {
		log.Fatalf("❌ Failed to parse templates: %v", err)
	}
	log.Println("✅ Templates loaded")
}

func setupRoutes(mux *http.ServeMux) {
	// Static assets
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("api/static")))
	mux.Handle("/static/", cacheControlMiddleware(staticHandler))

	// Health check
	mux.HandleFunc("/healthz", healthHandler)

	// Pages
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/dashboard", dashboardHandler)
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}

func gracefulShutdown(server *http.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	log.Println("🔽 Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Graceful shutdown failed: %v", err)
	}
	log.Println("✅ Server stopped cleanly")
}

/* -------------------- Handlers -------------------- */

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "base", map[string]interface{}{
		"Title":              "Home",
		"Year":               time.Now().Year(),
		"Description":        "Anonymous LLCs and Registered Agent services for solopreneurs who disrupt the status quo.",
		"IncludeDashboardJS": false,
	})
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "base", map[string]interface{}{
		"Title":              "Dashboard",
		"Year":               time.Now().Year(),
		"IncludeDashboardJS": true,
	})
}

/* -------------------- Middleware -------------------- */

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable") // 7 days
		next.ServeHTTP(w, r)
	})
}

/* -------------------- Utility -------------------- */

func renderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl, data); err != nil {
		log.Printf("❌ Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
