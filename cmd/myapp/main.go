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

/*
   --------------------------------------------------
   Global Template Registry
   --------------------------------------------------
*/
var templates *template.Template

func main() {
	log.Println("🔄 Starting AumKeeper server...")

	// Load .env if present (safe in prod)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found (ok in prod)")
	}

	// Graceful shutdown context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleSignals(cancel)

	// Load all templates once
	loadTemplates()

	// Router
	mux := http.NewServeMux()
	setupRoutes(mux)

	// Server
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
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	<-ctx.Done()
	gracefulShutdown(server, 10*time.Second)
}

/*
   --------------------------------------------------
   Template Loading
   --------------------------------------------------
*/

func loadTemplates() {
	var err error

	templates, err = template.ParseGlob("api/templates/*.html")
	if err != nil {
		log.Fatalf("❌ Template parse failure: %v", err)
	}

	log.Println("✅ Templates loaded")
}

/*
   --------------------------------------------------
   Routing
   --------------------------------------------------
*/

func setupRoutes(mux *http.ServeMux) {
	// Static files
	static := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("api/static")),
	)
	mux.Handle("/static/", cacheControlMiddleware(static))

	// Health
	mux.HandleFunc("/healthz", healthHandler)

	// Pages
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/dashboard", dashboardHandler)
}

/*
   --------------------------------------------------
   Handlers
   --------------------------------------------------
*/

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", map[string]interface{}{
		"Title":              "Home",
		"Description":        "Anonymous LLCs and Registered Agent services for modern founders.",
		"Year":               time.Now().Year(),
		"IncludeDashboardJS": false,
	})
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "dashboard", map[string]interface{}{
		"Title":              "Dashboard",
		"Year":               time.Now().Year(),
		"IncludeDashboardJS": true,
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

/*
   --------------------------------------------------
   Utilities
   --------------------------------------------------
*/

func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) {
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("❌ Template render error (%s): %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}

/*
   --------------------------------------------------
   Graceful Shutdown
   --------------------------------------------------
*/

func handleSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	log.Printf("🛑 Shutdown signal: %v", sig)
	cancel()
}

func gracefulShutdown(server *http.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Println("🔽 Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Shutdown failed: %v", err)
	}
	log.Println("✅ Server stopped cleanly")
}

/*
   --------------------------------------------------
   Middleware
   --------------------------------------------------
*/

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		next.ServeHTTP(w, r)
	})
}
