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

// Parse all templates once
var templates *template.Template

func main() {
	fmt.Println("🔄 Starting AumKeeper server...")

	// Load .env if present
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found or failed to load")
	}

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// OS signal handling
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		log.Printf("🛑 Received signal: %v", sig)
		cancel()
	}()

	// Load templates
	var err error
	templates, err = template.ParseGlob("api/templates/*.html")
	if err != nil {
		log.Fatalf("❌ Failed to load templates: %v", err)
	}
	log.Println("✅ Templates loaded")

	// Router
	mux := http.NewServeMux()

	// Static assets
	staticHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("api/static")),
	)
	mux.Handle("/static/", cacheControlMiddleware(staticHandler))

	// Health check
	mux.HandleFunc("/healthz", healthHandler)

	// Routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/dashboard", dashboardHandler)

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("🌐 Listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	log.Println("🔽 Shutting down server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
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
	data := map[string]interface{}{
		"Title":              "Home",
		"Year":               time.Now().Year(),
		"Description":        "Anonymous LLCs and Registered Agent services for solopreneurs who disrupt the status quo.",
		"IncludeDashboardJS": false,
	}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":              "Dashboard",
		"Year":               time.Now().Year(),
		"IncludeDashboardJS": true,
	}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* ---------------- Middleware ---------------- */

func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Cache-Control",
			"public, max-age=604800, immutable", // 7 days
		)
		next.ServeHTTP(w, r)
	})
}
