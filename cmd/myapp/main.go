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

	// Load .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found (ok in prod)")
	}

	// Graceful shutdown context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleSignals(cancel)

	// Load templates
	loadTemplates()

	// Router
	mux := http.NewServeMux()
	setupRoutes(mux)

	// HTTP Server
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

	// Parse base template first
	templates, err = template.New("").ParseGlob("api/templates/base.html")
	if err != nil {
		log.Fatalf("❌ Base template parse failure: %v", err)
	}

	// Parse all other page templates
	pageTemplates := []string{
		"api/templates/home.html",
		"api/templates/dashboard.html",
		"api/templates/staffs.html", // added staffs page
	}

	for _, tpl := range pageTemplates {
		templates, err = templates.ParseGlob(tpl)
		if err != nil {
			log.Fatalf("❌ Template parse failure (%s): %v", tpl, err)
		}
	}

	log.Println("✅ Templates loaded successfully")
}

/*
--------------------------------------------------
Routing
--------------------------------------------------
*/
func setupRoutes(mux *http.ServeMux) {
	// Static files
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("api/static")))
	mux.Handle("/static/", cacheControlMiddleware(static))

	// Health check
	mux.HandleFunc("/healthz", healthHandler)

	// Pages
	mux.HandleFunc("/staffs", staffsHandler)    // ✅ STAFFS page first
	mux.HandleFunc("/dashboard", dashboardHandler)
	mux.HandleFunc("/", homeHandler)            // root last
}

/*
--------------------------------------------------
Handlers
--------------------------------------------------
*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "home", map[string]any{
		"Title":       "AumKeeper — Modern ERP for Growing Businesses",
		"Description": "AumKeeper is a modern all-in-one SaaS ERP platform for SMBs.",
		"Year":        time.Now().Year(),
	})
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "dashboard", map[string]any{
		"Title":              "Dashboard",
		"Year":               time.Now().Year(),
		"IncludeDashboardJS": true,
	})
}

// STAFFS handler
func staffsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "staffs", map[string]any{
		"Title":          "Staff Manager",
		"Year":           time.Now().Year(),
		"IncludeStaffJS": true, // load staffs CSS & JS conditionally
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
func renderTemplate(w http.ResponseWriter, name string, data map[string]any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
